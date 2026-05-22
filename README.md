# jobsearch (working name: Pursuit)

Job search management for candidates. Closed beta. See `CLAUDE.md` for vision and roadmap, `TODO.md` for this week's items.

## Stack

- Go 1.24 backend (stdlib net/http, pgx for Postgres)
- Postgres 16
- Static frontend (HTML/CSS/vanilla JS)
- systemd + nginx on a Hetzner VM
- GitHub Actions for CI and deploy

## Local dev

```bash
# 1. start Postgres locally (or use docker run -e POSTGRES_PASSWORD=dev -p 5432:5432 postgres:16)
createdb jobsearch_dev

# 2. copy env
cp .env.example .env
# edit DATABASE_URL if needed

# 3. run
go run ./cmd/server
# server boots, applies migrations from ./migrations on startup
# magic-link emails are printed to stdout in dev (MAIL_DRIVER=log)
```

Then open http://localhost:8080.

## Layout

```
cmd/server/        entrypoint
internal/config/   env-driven config
internal/db/       pgx pool + migration runner
internal/auth/     sessions + magic-link tokens
internal/http/     server, middleware, handlers
internal/mail/     mail driver (log in dev, smtp/postmark later)
internal/llm/      Anthropic client (stub in v0.1)
migrations/        numbered SQL files, applied in order on boot
web/static/        frontend
deploy/            systemd unit + nginx sample + deploy notes
.github/workflows/ ci.yml, deploy.yml
```

## Deploy

See `deploy/README.md`. Short version: build a static Go binary in GH Actions, rsync to the Hetzner VM via SSH, `systemctl restart jobsearch`.
