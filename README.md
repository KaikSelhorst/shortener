# Shortener

A self-hosted link shortener with click analytics and webhooks.

Create short links per project, track clicks (device, browser, referrer, unique visitors), and get notified via webhooks when links are created, updated, deleted, or clicked.

## Stack

- **`api/`** — Go backend (standard library `net/http`, PostgreSQL)
- **`web/`** — SvelteKit frontend (Svelte 5, Tailwind CSS, Bun)
