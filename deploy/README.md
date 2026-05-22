# Deploy notes

Target: a single Hetzner Cloud VM (CX22 is plenty for closed beta).

## One-time server setup

```bash
# As root on a fresh Ubuntu 24.04 VM:
adduser --system --group --home /opt/jobsearch jobsearch
apt update && apt install -y postgresql nginx certbot python3-certbot-nginx rsync

# Postgres
sudo -u postgres createuser jobsearch
sudo -u postgres createdb -O jobsearch jobsearch
sudo -u postgres psql -c "ALTER USER jobsearch WITH PASSWORD '...'"  # use a real one

# App layout
install -d -o jobsearch -g jobsearch /opt/jobsearch/bin /opt/jobsearch/migrations /opt/jobsearch/web

# Drop in .env (copy from .env.example, fill DATABASE_URL, BASE_URL, ALLOWED_EMAILS)
$EDITOR /opt/jobsearch/.env
chown jobsearch:jobsearch /opt/jobsearch/.env
chmod 640 /opt/jobsearch/.env

# systemd
cp deploy/jobsearch.service /etc/systemd/system/jobsearch.service
systemctl daemon-reload
systemctl enable jobsearch

# nginx
cp deploy/nginx.conf /etc/nginx/sites-available/pursuit.conf
ln -s /etc/nginx/sites-available/pursuit.conf /etc/nginx/sites-enabled/
# Edit server_name, then:
certbot --nginx -d pursuit.example.com
nginx -t && systemctl reload nginx
```

## Deploys

Triggered by GitHub Actions on push to `main` (see `.github/workflows/deploy.yml`).

The workflow:
1. Builds a static Linux binary (`CGO_ENABLED=0 go build`).
2. Rsyncs `bin/server`, `migrations/`, and `web/static/` to the VM.
3. `systemctl restart jobsearch`.

Required GH Actions secrets:
- `DEPLOY_HOST` — `pursuit.example.com`
- `DEPLOY_USER` — usually `root` (or a deploy user with sudo for the restart)
- `DEPLOY_SSH_KEY` — private key authorized on the VM

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
