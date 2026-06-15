package webhook

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/KaikSelhorst/shortener/internal/repository"
	"github.com/KaikSelhorst/shortener/internal/service"
	"go.uber.org/zap"
)

func Start(ctx context.Context, repo repository.WebhookRepo, svc *service.WebhookService, client *http.Client, logger *zap.SugaredLogger) {
	if err := repo.ResetStuckDeliveries(ctx); err != nil {
		logger.Errorw("webhook: reset stuck deliveries", "error", err)
	}
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			process(ctx, repo, svc, client, logger)
		}
	}
}

func process(ctx context.Context, repo repository.WebhookRepo, svc *service.WebhookService, client *http.Client, logger *zap.SugaredLogger) {
	pending, err := repo.ClaimPendingDeliveries(ctx, 10)
	if err != nil {
		logger.Errorw("webhook: claim deliveries", "error", err)
		return
	}
	for _, p := range pending {
		go deliver(ctx, repo, svc, client, logger, p)
	}
}

func deliver(ctx context.Context, repo repository.WebhookRepo, svc *service.WebhookService, client *http.Client, logger *zap.SugaredLogger, p *repository.PendingDelivery) {
	secret, err := svc.DecryptSecret(p.EncryptedSecret)
	if err != nil {
		logger.Errorw("webhook: decrypt secret", "id", p.ID, "error", err)
		markFailed(ctx, repo, p.ID)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.TargetURL, bytes.NewReader(p.Payload))
	if err != nil {
		logger.Errorw("webhook: build request", "id", p.ID, "error", err)
		markFailed(ctx, repo, p.ID)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Event", p.Event)
	req.Header.Set("X-Webhook-Signature", service.Sign(p.Payload, secret))

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var code *int
		if resp != nil {
			c := resp.StatusCode
			code = &c
		}
		scheduleRetry(ctx, repo, logger, p.ID, p.Attempts, code)
		return
	}

	code := resp.StatusCode
	if err := repo.UpdateDeliveryStatus(ctx, p.ID, "delivered", &code, nil); err != nil {
		logger.Errorw("webhook: mark delivered", "id", p.ID, "error", err)
	}
}

func scheduleRetry(ctx context.Context, repo repository.WebhookRepo, logger *zap.SugaredLogger, id string, attempts int, responseStatus *int) {
	delay := retryDelay(attempts)
	if delay == 0 {
		markFailed(ctx, repo, id)
		return
	}
	nextRetry := time.Now().Add(delay)
	if err := repo.UpdateDeliveryStatus(ctx, id, "pending", responseStatus, &nextRetry); err != nil {
		logger.Errorw("webhook: schedule retry", "id", id, "error", err)
	}
}

func markFailed(ctx context.Context, repo repository.WebhookRepo, id string) {
	_ = repo.UpdateDeliveryStatus(ctx, id, "failed", nil, nil)
}

func retryDelay(attempts int) time.Duration {
	switch attempts {
	case 1:
		return 30 * time.Second
	case 2:
		return 5 * time.Minute
	case 3:
		return 30 * time.Minute
	default:
		return 0
	}
}
