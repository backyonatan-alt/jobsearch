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
# 1. start Postgres locally
createdb jobsearch_dev

# 2. copy env, fill DATABASE_URL + GOOGLE_CLIENT_ID + GOOGLE_CLIENT_SECRET
cp .env.example .env

# 3. build the frontend once (Vite outputs to web/build/)
cd web && pnpm install && pnpm build && cd ..

# 4. run the Go server
go run ./cmd/server

# server boots, applies migrations from ./migrations, serves web/build/
```

Then open http://localhost:8080.

For frontend iteration, run Vite's dev server in a second terminal:

```bash
cd web && pnpm dev
```

Vite serves on :5173 with HMR and proxies /api, /auth, /healthz to the Go server on :8080.

## Layout

```
cmd/server/        Go entrypoint
internal/config/   env-driven config
internal/db/       pgx pool + migration runner
internal/auth/     sessions + Google OAuth
internal/httpsrv/  server, middleware, handlers
migrations/        numbered SQL files, applied in order on boot
web/               SvelteKit frontend (Vite + adapter-static)
  src/             pages, components, $lib helpers
  static/          static assets passed through to build/
  build/           Vite output (gitignored, served by Go)
deploy/            bootstrap.sh, systemd unit, nginx sample
.github/workflows/ ci.yml, deploy.yml
```

## Deploy

See `deploy/README.md`. Short version: build a static Go binary in GH Actions, rsync to the Hetzner VM via SSH, `systemctl restart jobsearch`.
