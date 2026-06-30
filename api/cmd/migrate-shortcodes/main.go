// migrate-shortcodes regenerates every auto-generated short_code using the
// new base-62 + Knuth scramble algorithm, replacing codes produced by the
// old sqids-go library.
//
// Custom codes (manually set by users) cannot be distinguished from generated
// ones in the database, so this command regenerates ALL codes. Run it once
// before or immediately after deploying the new service binary.
//
// Usage:
//
//	DATABASE_URL=postgres://... SHORTCODE_SECRET=... go run ./cmd/migrate-shortcodes
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/KaikSelhorst/shortener/internal/database"
	"github.com/KaikSelhorst/shortener/internal/service"
	"github.com/jackc/pgx/v5"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	secret := os.Getenv("SHORTCODE_SECRET")
	if secret == "" {
		log.Fatal("SHORTCODE_SECRET is required")
	}

	svc, err := service.NewShortcodeService([]byte(secret))
	if err != nil {
		log.Fatalf("shortcode service: %v", err)
	}

	ctx := context.Background()
	db, err := database.NewDatabase(ctx, dbURL)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer db.Close()

	rows, err := db.Pool.Query(ctx, `SELECT id, short_code FROM links ORDER BY id`)
	if err != nil {
		log.Fatalf("query links: %v", err)
	}

	type row struct {
		id        int64
		shortCode string
	}
	var links []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.id, &r.shortCode); err != nil {
			log.Fatalf("scan: %v", err)
		}
		links = append(links, r)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		log.Fatalf("rows: %v", err)
	}

	fmt.Printf("found %d links to migrate\n\n", len(links))

	batch := &pgx.Batch{}
	skipped := 0
	for _, l := range links {
		newCode, err := svc.GenerateShortCode(uint64(l.id))
		if err != nil {
			log.Fatalf("generate shortcode for id=%d: %v", l.id, err)
		}
		if newCode == l.shortCode {
			fmt.Printf("· %d  %s  (already correct)\n", l.id, l.shortCode)
			skipped++
			continue
		}
		fmt.Printf("✓ %d  %s  →  %s\n", l.id, l.shortCode, newCode)
		batch.Queue(`UPDATE links SET short_code = $1 WHERE id = $2`, newCode, l.id)
	}

	if batch.Len() == 0 {
		fmt.Printf("\nall %d links already up to date\n", skipped)
		return
	}

	br := db.Pool.SendBatch(ctx, batch)
	for range batch.Len() {
		if _, err := br.Exec(); err != nil {
			_ = br.Close()
			log.Fatalf("update batch: %v", err)
		}
	}
	if err := br.Close(); err != nil {
		log.Fatalf("close batch: %v", err)
	}

	fmt.Printf("\nmigrated %d links, skipped %d\n", batch.Len(), skipped)
}
