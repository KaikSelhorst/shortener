package service

import (
	"context"
	"sync"
	"time"

	"github.com/KaikSelhorst/shortener/internal/model"
	"github.com/KaikSelhorst/shortener/internal/repository"
	"go.uber.org/zap"
)

const (
	trackerBufSize       = 512
	trackerBatchSize     = 100
	trackerFlushInterval = 5 * time.Second
)

type TrackerService struct {
	repo   *repository.ClickRepository
	logger *zap.SugaredLogger
	queue  chan model.Click
	done   chan struct{}
	wg     sync.WaitGroup
}

func NewTrackerService(repo *repository.ClickRepository, logger *zap.SugaredLogger) *TrackerService {
	t := &TrackerService{
		repo:   repo,
		logger: logger,
		queue:  make(chan model.Click, trackerBufSize),
		done:   make(chan struct{}),
	}
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		t.run()
	}()
	return t
}

func (t *TrackerService) Record(click model.Click) {
	select {
	case t.queue <- click:
	default:
		t.logger.Warn("tracker queue full, dropping click")
	}
}

func (t *TrackerService) Shutdown() {
	close(t.done)
	t.wg.Wait()
}

func (t *TrackerService) run() {
	ticker := time.NewTicker(trackerFlushInterval)
	defer ticker.Stop()

	batch := make([]model.Click, 0, trackerBatchSize)

	flush := func() {
		if len(batch) == 0 {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := t.repo.BatchInsert(ctx, batch); err != nil {
			t.logger.Errorw("tracker batch insert failed", "error", err)
		}
		batch = batch[:0]
	}

	for {
		select {
		case click := <-t.queue:
			batch = append(batch, click)
			if len(batch) >= trackerBatchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		case <-t.done:
			for {
				select {
				case click := <-t.queue:
					batch = append(batch, click)
				default:
					flush()
					return
				}
			}
		}
	}
}
