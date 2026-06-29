// Seed populates the database with demo data for local development.
// It is idempotent: running it multiple times is safe.
//
// Usage:
//
//	DATABASE_URL=postgres://... IP_HASH_SECRET=... go run ./cmd/seed
//	make seed
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"time"

	"github.com/KaikSelhorst/shortener/internal/database"
	"github.com/KaikSelhorst/shortener/internal/service"
	"github.com/KaikSelhorst/shortener/migrations"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const (
	demoEmail    = "demo@shortener.dev"
	demoPassword = "demo1234"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	ipSecret := os.Getenv("IP_HASH_SECRET")
	if ipSecret == "" {
		log.Fatal("IP_HASH_SECRET is required")
	}
	webhookKey := os.Getenv("WEBHOOK_SECRET_KEY")
	if webhookKey == "" {
		log.Fatal("WEBHOOK_SECRET_KEY is required")
	}

	if err := migrations.Run(dbURL); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	ctx := context.Background()
	db, err := database.NewDatabase(ctx, dbURL)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer db.Close()

	shortcodes, err := service.NewShortcodeService()
	if err != nil {
		log.Fatalf("shortcode service: %v", err)
	}

	webhookSvc := service.NewWebhookService(nil, []byte(webhookKey))

	userID := ensureUser(ctx, db.Pool)
	projectIDs := ensureProjects(ctx, db.Pool, userID)
	links := ensureLinks(ctx, db.Pool, shortcodes, projectIDs)
	ensureClicks(ctx, db.Pool, links, ipSecret)
	webhookIDs := ensureWebhooks(ctx, db.Pool, projectIDs, webhookSvc)
	ensureWebhookDeliveries(ctx, db.Pool, webhookIDs)

	fmt.Println()
	fmt.Println("  Ready. Start the server and log in:")
	fmt.Printf("    Email:    %s\n", demoEmail)
	fmt.Printf("    Password: %s\n", demoPassword)
}

// --- User ----------------------------------------------------------------

func ensureUser(ctx context.Context, pool *pgxpool.Pool) int64 {
	var id int64
	err := pool.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, demoEmail).Scan(&id)
	if err == nil {
		fmt.Printf("· user        %s\n", demoEmail)
		return id
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(demoPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("bcrypt: %v", err)
	}
	if err = pool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`,
		demoEmail, string(hash),
	).Scan(&id); err != nil {
		log.Fatalf("insert user: %v", err)
	}
	fmt.Printf("✓ user        %s\n", demoEmail)
	return id
}

// --- Projects ------------------------------------------------------------

var projectSeeds = []struct{ name, slug string }{
	{"Marketing", "marketing"},
	{"Tech Blog", "tech-blog"},
	{"Social Media", "social-media"},
}

func ensureProjects(ctx context.Context, pool *pgxpool.Pool, userID int64) map[string]int64 {
	ids := make(map[string]int64, len(projectSeeds))
	for _, p := range projectSeeds {
		var id int64
		err := pool.QueryRow(ctx,
			`SELECT id FROM projects WHERE user_id = $1 AND slug = $2`, userID, p.slug,
		).Scan(&id)
		if err != nil {
			if err = pool.QueryRow(ctx,
				`INSERT INTO projects (user_id, name, slug) VALUES ($1, $2, $3) RETURNING id`,
				userID, p.name, p.slug,
			).Scan(&id); err != nil {
				log.Fatalf("insert project %q: %v", p.name, err)
			}
			fmt.Printf("✓ project     %s\n", p.name)
		} else {
			fmt.Printf("· project     %s\n", p.name)
		}
		ids[p.slug] = id
	}
	return ids
}

// --- Links ---------------------------------------------------------------

var linkSeeds = []struct {
	project    string
	url        string
	title      string
	customCode string // empty = auto-generate
	maxClicks  int64  // 0 = unlimited
}{
	{"marketing", "https://github.com/KaikSelhorst/shortener", "Shortener on GitHub", "repo", 0},
	{"marketing", "https://example.com/landing", "Landing Page", "landing", 500},
	{"marketing", "https://example.com/pricing", "Pricing", "pricing", 0},
	{"marketing", "https://example.com/docs", "Documentation", "", 0},
	{"tech-blog", "https://go.dev/blog/intro-generics", "Intro to Go Generics", "", 0},
	{"tech-blog", "https://svelte.dev/blog/runes", "Svelte Runes", "svelte-runes", 100},
	{"tech-blog", "https://www.postgresql.org/docs/current/", "PostgreSQL Docs", "pg-docs", 0},
	{"tech-blog", "https://tailwindcss.com/blog/tailwindcss-v4", "Tailwind CSS v4", "", 0},
	{"social-media", "https://twitter.com/golang", "Go on Twitter/X", "golang-x", 0},
	{"social-media", "https://linkedin.com/company/postgresql", "PostgreSQL on LinkedIn", "", 50},
	{"social-media", "https://youtube.com/@ThePrimeagen", "ThePrimeagen on YouTube", "prime", 0},
}

type seededLink struct {
	id        int64
	shortCode string
}

func ensureLinks(ctx context.Context, pool *pgxpool.Pool, svc *service.ShortcodeService, projectIDs map[string]int64) []seededLink {
	result := make([]seededLink, 0, len(linkSeeds))
	for _, l := range linkSeeds {
		pid := projectIDs[l.project]
		var id int64
		var code string
		err := pool.QueryRow(ctx,
			`SELECT id, COALESCE(short_code, '') FROM links WHERE project_id = $1 AND original_url = $2`,
			pid, l.url,
		).Scan(&id, &code)
		if err != nil {
			var maxClicks *int64
			if l.maxClicks > 0 {
				mc := l.maxClicks
				maxClicks = &mc
			}
			if l.customCode != "" {
				if err = pool.QueryRow(ctx,
					`INSERT INTO links (project_id, original_url, title, short_code, max_clicks) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
					pid, l.url, l.title, l.customCode, maxClicks,
				).Scan(&id); err != nil {
					log.Fatalf("insert link %q: %v", l.url, err)
				}
				code = l.customCode
			} else {
				if err = pool.QueryRow(ctx,
					`INSERT INTO links (project_id, original_url, title, max_clicks) VALUES ($1, $2, $3, $4) RETURNING id`,
					pid, l.url, l.title, maxClicks,
				).Scan(&id); err != nil {
					log.Fatalf("insert link %q: %v", l.url, err)
				}
				code, err = svc.GenerateShortCode(uint64(id))
				if err != nil {
					log.Fatalf("generate shortcode for id=%d: %v", id, err)
				}
				if _, err = pool.Exec(ctx,
					`UPDATE links SET short_code = $1 WHERE id = $2`, code, id,
				); err != nil {
					log.Fatalf("update short_code: %v", err)
				}
			}
			fmt.Printf("✓ link        /%s  →  %s\n", code, l.url)
		} else {
			fmt.Printf("· link        /%s\n", code)
		}
		result = append(result, seededLink{id: id, shortCode: code})
	}
	return result
}

// --- Clicks --------------------------------------------------------------

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:125.0) Gecko/20100101 Firefox/125.0",
	"Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0",
	"Mozilla/5.0 (Linux; Android 12; SM-G998B) AppleWebKit/537.36 SamsungBrowser/19.0 Chrome/102.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64; rv:125.0) Gecko/20100101 Firefox/125.0",
	"Mozilla/5.0 (iPad; CPU OS 17_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
}

var referers = []string{
	"", "", "", // direct (higher weight)
	"https://www.google.com/search?q=url+shortener",
	"https://www.google.com/search?q=short+links+golang",
	"https://www.instagram.com/",
	"https://l.instagram.com/?u=https%3A%2F%2Fshortener.dev",
	"https://twitter.com/home",
	"https://t.co/abc123xyz",
	"https://www.facebook.com/sharer",
	"https://discord.com/channels/123456789/987654321",
	"https://discord.gg/golang",
	"https://www.linkedin.com/feed/",
	"https://www.youtube.com/watch?v=dQw4w9WgXcQ",
	"https://news.ycombinator.com/item?id=12345",
}

var sampleIPs = []string{
	"203.0.113.10", "203.0.113.42", "198.51.100.20",
	"198.51.100.77", "192.0.2.15", "192.0.2.88",
	"185.220.101.1", "104.16.100.1", "45.33.32.156",
	"151.101.1.195", "8.8.8.8", "1.1.1.1",
}

// --- Webhooks ------------------------------------------------------------

var webhookSeeds = []struct {
	project string
	url     string
	events  []string
}{
	{"marketing", "https://hooks.example.com/marketing", []string{"link.created", "link.deleted"}},
	{"tech-blog", "https://hooks.example.com/tech-blog", []string{"link.clicked", "link.created", "link.updated", "link.deleted"}},
	{"social-media", "https://hooks.example.com/social", []string{"link.clicked"}},
}

func ensureWebhooks(ctx context.Context, pool *pgxpool.Pool, projectIDs map[string]int64, svc *service.WebhookService) map[string]string {
	ids := make(map[string]string, len(webhookSeeds))
	for _, w := range webhookSeeds {
		pid := projectIDs[w.project]
		var id string
		err := pool.QueryRow(ctx,
			`SELECT id FROM webhooks WHERE project_id = $1 AND url = $2`, pid, w.url,
		).Scan(&id)
		if err != nil {
			plain, err := service.GenerateWebhookSecret()
			if err != nil {
				log.Fatalf("generate webhook secret: %v", err)
			}
			encrypted, err := svc.EncryptSecret(plain)
			if err != nil {
				log.Fatalf("encrypt webhook secret: %v", err)
			}
			if err = pool.QueryRow(ctx,
				`INSERT INTO webhooks (project_id, url, secret, events) VALUES ($1, $2, $3, $4) RETURNING id`,
				pid, w.url, encrypted, w.events,
			).Scan(&id); err != nil {
				log.Fatalf("insert webhook %q: %v", w.url, err)
			}
			fmt.Printf("✓ webhook     %s  [%s]\n", w.url, w.project)
		} else {
			fmt.Printf("· webhook     %s\n", w.url)
		}
		ids[w.project] = id
	}
	return ids
}

// --- Webhook deliveries --------------------------------------------------

func ensureWebhookDeliveries(ctx context.Context, pool *pgxpool.Pool, webhookIDs map[string]string) {
	type deliverySeed struct {
		project        string
		event          string
		status         string
		attempts       int
		responseStatus *int
		daysAgo        int
	}

	ok200 := 200
	err500 := 500
	err404 := 404

	seeds := []deliverySeed{
		{"marketing", "link.created", "delivered", 1, &ok200, 1},
		{"marketing", "link.created", "delivered", 1, &ok200, 3},
		{"marketing", "link.deleted", "delivered", 1, &ok200, 5},
		{"marketing", "link.deleted", "failed", 4, &err500, 7},
		{"tech-blog", "link.clicked", "delivered", 1, &ok200, 0},
		{"tech-blog", "link.clicked", "delivered", 1, &ok200, 1},
		{"tech-blog", "link.clicked", "delivered", 2, &ok200, 2},
		{"tech-blog", "link.created", "delivered", 1, &ok200, 3},
		{"tech-blog", "link.updated", "failed", 4, &err404, 4},
		{"tech-blog", "link.deleted", "pending", 1, nil, 0},
		{"social-media", "link.clicked", "delivered", 1, &ok200, 0},
		{"social-media", "link.clicked", "delivered", 1, &ok200, 2},
		{"social-media", "link.clicked", "failed", 4, &err500, 6},
	}

	for project, webhookID := range webhookIDs {
		var count int64
		_ = pool.QueryRow(ctx,
			`SELECT COUNT(*) FROM webhook_deliveries WHERE webhook_id = $1`, webhookID,
		).Scan(&count)
		if count > 0 {
			fmt.Printf("· deliveries  webhook#%s already present\n", webhookID)
			continue
		}

		inserted := 0
		for _, s := range seeds {
			if s.project != project {
				continue
			}
			payload, _ := json.Marshal(map[string]any{
				"event":      s.event,
				"project_id": webhookIDs[s.project],
			})
			createdAt := time.Now().AddDate(0, 0, -s.daysAgo)
			if _, err := pool.Exec(ctx,
				`INSERT INTO webhook_deliveries (webhook_id, event, payload, status, attempts, response_status, created_at)
				 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
				webhookID, s.event, payload, s.status, s.attempts, s.responseStatus, createdAt,
			); err != nil {
				log.Fatalf("insert delivery for webhook#%s: %v", webhookID, err)
			}
			inserted++
		}
		fmt.Printf("✓ deliveries  %d inserted for webhook#%s\n", inserted, webhookID)
	}
}

func ensureClicks(ctx context.Context, pool *pgxpool.Pool, links []seededLink, ipHashSecret string) {
	ids := make([]int64, len(links))
	for i, l := range links {
		ids[i] = l.id
	}

	var count int64
	_ = pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM clicks WHERE link_id = ANY($1)`, ids,
	).Scan(&count)

	if count > 0 {
		fmt.Printf("· clicks      %d already present\n", count)
		return
	}

	now := time.Now()
	rows := make([][]any, 0, len(links)*30)

	for _, link := range links {
		n := 20 + rand.IntN(20) // 20–39 clicks per link
		for range n {
			daysAgo := rand.IntN(30)
			hoursAgo := rand.IntN(24)
			ts := now.AddDate(0, 0, -daysAgo).Add(-time.Duration(hoursAgo) * time.Hour)

			ua := userAgents[rand.IntN(len(userAgents))]
			ref := referers[rand.IntN(len(referers))]
			rawIP := sampleIPs[rand.IntN(len(sampleIPs))]

			ipHash := service.HashIP(rawIP, ipHashSecret)
			deviceType := service.ParseDeviceType(ua)
			refSource := service.ParseReferrerSource(ref)
			browser := service.ParseBrowserName(ua)

			var refVal any
			if ref != "" {
				refVal = ref
			}

			rows = append(rows, []any{
				link.id, ua, ipHash, refVal,
				deviceType, refSource, browser, ts,
			})
		}
	}

	_, err := pool.CopyFrom(ctx,
		pgx.Identifier{"clicks"},
		[]string{"link_id", "user_agent", "ip_hash", "referer", "device_type", "referrer_source", "browser", "created_at"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Fatalf("insert clicks: %v", err)
	}
	fmt.Printf("✓ clicks      %d inserted across %d links\n", len(rows), len(links))
}
