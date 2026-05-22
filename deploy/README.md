# Deploy notes

Target: a single Hetzner Cloud VM (CX22 is plenty for closed beta).

We use `<public-ip>.nip.io` as the hostname during the beta so we don't have to
buy a domain. nip.io resolves any `1.2.3.4.nip.io` to `1.2.3.4`, and Let's
Encrypt happily issues real certs against it. Swap to a real domain later by
changing `BASE_URL` in `.env`, the nginx `server_name`, and rerunning certbot.

## One-time server setup

```bash
# 1. Create a fresh Ubuntu 24.04 VM in the Hetzner console (CX22 is enough).
#    Note the public IP.

# 2. SSH in as root and run the bootstrap script:
ssh root@<ip> bash -c "'curl -fsSL https://raw.githubusercontent.com/backyonatan-alt/jobsearch/main/deploy/bootstrap.sh | ADMIN_EMAIL=you@example.com bash'"

# Or, if you're on a feature branch and the script isn't on main yet:
scp deploy/bootstrap.sh root@<ip>:/root/
ssh root@<ip> "ADMIN_EMAIL=you@example.com bash /root/bootstrap.sh"
```

The script installs Postgres + nginx + certbot, creates the `jobsearch` user
and DB, writes `/opt/jobsearch/.env`, installs the systemd unit, configures
nginx for `<ip>.nip.io`, and runs certbot to issue a TLS cert. It's
re-runnable — anything that already exists is left alone.

## Deploys

Triggered by GitHub Actions on push to `main` (see `.github/workflows/deploy.yml`).

The workflow:
1. Builds a static Linux binary (`CGO_ENABLED=0 go build`).
2. Rsyncs `bin/server`, `migrations/`, and `web/static/` to the VM.
3. `systemctl restart jobsearch`.

Required GH Actions secrets:
- `DEPLOY_HOST` — e.g. `1.2.3.4.nip.io`
- `DEPLOY_USER` — usually `root` (or a deploy user with sudo for the restart)
- `DEPLOY_SSH_KEY` — private key authorized on the VM

Required GH Actions variable (not secret):
- `DEPLOY_ENABLED=true` — gates the deploy workflow until you're ready

Manual deploy from a laptop is also fine:

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/server ./cmd/server
rsync -avz --delete bin/ migrations/ web/ user@host:/opt/jobsearch/
ssh user@host 'sudo systemctl restart jobsearch && sleep 1 && curl -fsS http://localhost:8080/healthz'
```

## Verify after deploy

```bash
ssh user@host 'systemctl status jobsearch --no-pager | head -20'
curl -fsS https://pursuit.example.com/healthz
```
