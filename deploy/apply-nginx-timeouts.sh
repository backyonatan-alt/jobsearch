#!/usr/bin/env bash
# Idempotently add AI-endpoint proxy timeouts to the live nginx config.
# Safe to run on every deploy — checks if the new location is already
# present and bails out. Preserves whatever certbot wrote (SSL listeners,
# cert paths, etc.) by inserting *before* the catch-all `location / {`.

set -euo pipefail

CONF=/etc/nginx/sites-available/pursuit.conf

if [ ! -f "$CONF" ]; then
  echo "nginx config not found at $CONF — skipping (fresh box?)"
  exit 0
fi

if grep -q "dossier/refresh" "$CONF"; then
  echo "nginx config already has AI proxy timeouts — skipping"
  exit 0
fi

echo "patching nginx config — adding AI proxy timeouts"
cp "$CONF" "$CONF.bak.$(date +%s)"

awk '
  /^[[:space:]]*location \/ \{/ {
    print "    location ~ ^/api/applications/[0-9]+/(dossier/refresh|interviews/parse)$ {";
    print "        proxy_pass         http://127.0.0.1:8080;";
    print "        proxy_http_version 1.1;";
    print "        proxy_set_header   Host              $host;";
    print "        proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;";
    print "        proxy_set_header   X-Forwarded-Proto $scheme;";
    print "        proxy_read_timeout 180s;";
    print "        proxy_connect_timeout 5s;";
    print "    }";
    print "    location = /api/applications/parse {";
    print "        proxy_pass         http://127.0.0.1:8080;";
    print "        proxy_http_version 1.1;";
    print "        proxy_set_header   Host              $host;";
    print "        proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;";
    print "        proxy_set_header   X-Forwarded-Proto $scheme;";
    print "        proxy_read_timeout 90s;";
    print "        proxy_connect_timeout 5s;";
    print "    }";
    print "";
  }
  { print }
' "$CONF" > /tmp/pursuit.conf.new

mv /tmp/pursuit.conf.new "$CONF"
nginx -t
systemctl reload nginx
echo "nginx patched and reloaded — AI endpoints now allow up to 180s"
