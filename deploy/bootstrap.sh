#!/usr/bin/env bash
# One-shot bootstrap for a fresh Ubuntu 24.04 Hetzner VM.
#
# What this does:
#   1. Installs Postgres, nginx, certbot, rsync.
#   2. Creates a `jobsearch` system user and /opt/jobsearch layout.
#   3. Creates the Postgres role + database with a generated password.
#   4. Drops a starter /opt/jobsearch/.env (you fill ALLOWED_EMAILS after).
#   5. Installs the systemd unit and nginx config (nip.io hostname).
#   6. Runs certbot to get a Let's Encrypt cert for <ip>.nip.io.
#   7. Prints the GH Actions secrets to set.
#
# Re-runnable: skips anything that already exists.
#
# Usage (as root on the VM):
#   curl -fsSL https://raw.githubusercontent.com/backyonatan-alt/jobsearch/main/deploy/bootstrap.sh | bash
# or:
#   scp deploy/bootstrap.sh root@<ip>:/root/
#   ssh root@<ip> bash /root/bootstrap.sh

set -euo pipefail

if [[ "$(id -u)" -ne 0 ]]; then
  echo "Run as root." >&2
  exit 1
fi

PUBLIC_IP="${PUBLIC_IP:-$(curl -fsS https://ipv4.icanhazip.com || true)}"
if [[ -z "${PUBLIC_IP}" ]]; then
  echo "Could not detect public IP; export PUBLIC_IP=… and re-run." >&2
  exit 1
fi
HOSTNAME_NIP="${HOSTNAME_NIP:-${PUBLIC_IP}.nip.io}"
: "${ADMIN_EMAIL:?set ADMIN_EMAIL=you@example.com (used for Lets Encrypt cert issuance)}"

echo "==> Public IP:     ${PUBLIC_IP}"
echo "==> Hostname:      ${HOSTNAME_NIP}"
echo "==> Admin email:   ${ADMIN_EMAIL}"
echo

echo "==> apt update + install packages"
export DEBIAN_FRONTEND=noninteractive
apt-get update -qq
apt-get install -y -qq postgresql nginx certbot python3-certbot-nginx rsync ufw curl

echo "==> Configure UFW (allow ssh + http + https)"
ufw --force reset >/dev/null
ufw default deny incoming >/dev/null
ufw default allow outgoing >/dev/null
ufw allow OpenSSH >/dev/null
ufw allow 'Nginx Full' >/dev/null
ufw --force enable >/dev/null

echo "==> Create jobsearch system user"
if ! id jobsearch >/dev/null 2>&1; then
  adduser --system --group --home /opt/jobsearch jobsearch
fi
install -d -o jobsearch -g jobsearch -m 750 /opt/jobsearch
install -d -o jobsearch -g jobsearch -m 750 /opt/jobsearch/bin
install -d -o jobsearch -g jobsearch -m 750 /opt/jobsearch/migrations
install -d -o jobsearch -g jobsearch -m 750 /opt/jobsearch/web

echo "==> Postgres role + database"
if ! sudo -u postgres psql -tAc "SELECT 1 FROM pg_roles WHERE rolname='jobsearch'" | grep -q 1; then
  PG_PASS="$(openssl rand -hex 24)"
  sudo -u postgres psql -c "CREATE USER jobsearch WITH PASSWORD '${PG_PASS}'"
  sudo -u postgres createdb -O jobsearch jobsearch
  echo "${PG_PASS}" > /opt/jobsearch/.pgpass
  chown jobsearch:jobsearch /opt/jobsearch/.pgpass
  chmod 600 /opt/jobsearch/.pgpass
else
  echo "    (jobsearch role already exists, leaving alone)"
  PG_PASS="$(cat /opt/jobsearch/.pgpass 2>/dev/null || echo 'UNKNOWN_SEE_PG')"
fi

echo "==> Write /opt/jobsearch/.env"
if [[ ! -f /opt/jobsearch/.env ]]; then
  cat > /opt/jobsearch/.env <<EOF
PORT=8080
BASE_URL=https://${HOSTNAME_NIP}
BRAND_NAME=Pursuit

DATABASE_URL=postgres://jobsearch:${PG_PASS}@localhost:5432/jobsearch?sslmode=disable

SESSION_COOKIE_NAME=pursuit_session
SESSION_TTL_HOURS=720

# Google OAuth. Fill these in after creating an OAuth Web client at
# https://console.cloud.google.com/apis/credentials with the authorized
# redirect URI:  https://${HOSTNAME_NIP}/auth/google/callback
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=

ANTHROPIC_API_KEY=

# Closed beta gate. Add comma-separated invitee emails before opening up.
ALLOWED_EMAILS=${ADMIN_EMAIL}
EOF
  chown jobsearch:jobsearch /opt/jobsearch/.env
  chmod 640 /opt/jobsearch/.env
else
  echo "    (.env already exists, leaving alone)"
fi

echo "==> Install systemd unit"
cat > /etc/systemd/system/jobsearch.service <<'EOF'
[Unit]
Description=Pursuit (jobsearch) HTTP server
After=network.target postgresql.service

[Service]
Type=simple
User=jobsearch
Group=jobsearch
WorkingDirectory=/opt/jobsearch
EnvironmentFile=/opt/jobsearch/.env
ExecStart=/opt/jobsearch/bin/server
Restart=always
RestartSec=2

NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
PrivateTmp=true
ReadWritePaths=/opt/jobsearch

[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload
systemctl enable jobsearch >/dev/null

echo "==> Install nginx site"
cat > /etc/nginx/sites-available/pursuit.conf <<EOF
server {
    listen 80;
    listen [::]:80;
    server_name ${HOSTNAME_NIP};

    client_max_body_size 25M;

    # Long-running AI endpoints (dossier with web_search, interview parse with
    # vision) can take 60-150s. Carve them out with a longer read timeout so
    # nginx doesn't 504 before Claude finishes.
    location ~ ^/api/applications/[0-9]+/(dossier/refresh|interviews/parse)\$ {
        proxy_pass         http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header   Host              \$host;
        proxy_set_header   X-Forwarded-For   \$proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto \$scheme;
        proxy_read_timeout 180s;
        proxy_connect_timeout 5s;
    }
    location = /api/applications/parse {
        proxy_pass         http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header   Host              \$host;
        proxy_set_header   X-Forwarded-For   \$proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto \$scheme;
        proxy_read_timeout 90s;
        proxy_connect_timeout 5s;
    }

    location / {
        proxy_pass         http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header   Host              \$host;
        proxy_set_header   X-Forwarded-For   \$proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto \$scheme;
        proxy_read_timeout 60s;
        proxy_next_upstream error timeout http_502 http_503 http_504;
        proxy_connect_timeout 2s;
    }
}
EOF
ln -sf /etc/nginx/sites-available/pursuit.conf /etc/nginx/sites-enabled/pursuit.conf
rm -f /etc/nginx/sites-enabled/default
nginx -t
systemctl reload nginx

echo "==> Issue Let's Encrypt cert for ${HOSTNAME_NIP}"
if ! [[ -f "/etc/letsencrypt/live/${HOSTNAME_NIP}/fullchain.pem" ]]; then
  certbot --nginx --non-interactive --agree-tos \
    --email "${ADMIN_EMAIL}" -d "${HOSTNAME_NIP}" --redirect
else
  echo "    (cert already issued, skipping)"
fi

echo
echo "============================================================"
echo " Bootstrap complete."
echo
echo " URL:        https://${HOSTNAME_NIP}"
echo " App user:   jobsearch"
echo " App dir:    /opt/jobsearch"
echo " systemd:    systemctl status jobsearch"
echo " logs:       journalctl -u jobsearch -f"
echo
echo " Next steps:"
echo "   1. Add a deploy SSH key:"
echo "        ssh-keygen -t ed25519 -f /root/.ssh/deploy_jobsearch -N ''"
echo "        cat /root/.ssh/deploy_jobsearch.pub >> /root/.ssh/authorized_keys"
echo "      Then copy the PRIVATE key (deploy_jobsearch) into the GH Actions secret"
echo "      DEPLOY_SSH_KEY."
echo
echo "   2. Set these GitHub Actions secrets on the repo:"
echo "        DEPLOY_HOST     = ${HOSTNAME_NIP}"
echo "        DEPLOY_USER     = root   (or a sudo-capable deploy user)"
echo "        DEPLOY_SSH_KEY  = <private key from step 1>"
echo "      And set this Actions VARIABLE (not secret):"
echo "        DEPLOY_ENABLED  = true"
echo
echo "   3. Push to main. The deploy workflow will rsync the binary and restart."
echo
echo "   4. Hit https://${HOSTNAME_NIP}, enter ${ADMIN_EMAIL}, then run:"
echo "        journalctl -u jobsearch -n 50 | grep -i 'magic link'"
echo "      to grab your sign-in link."
echo "============================================================"
