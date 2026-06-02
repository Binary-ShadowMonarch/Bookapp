#!/usr/bin/env bash
set -euo pipefail

# ─── CONFIG ────────────────────────────────────────────────────────────────────
API_HOST="localhost:8080"

# test user credentials
EMAIL="integration@example.com"
PASSWORD="secret123"

# Postgres connection (psql CLI must be installed)
export PGPASSWORD='pOnw14WvDHWhV90F7KAq4KUj'
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="admin"
DB_NAME="bookapp"

# SeaweedFS filer
SEAWEEDFS_FILER="localhost:8888"

# ─── 1) REGISTER VIA API ─────────────────────────────────────────────────────
echo "1) Registering user $EMAIL…"
curl -s -X POST \
  -d "mail=$EMAIL&password=$PASSWORD" \
  "http://${API_HOST}/register" \
  && echo " → OK" || { echo " → FAILED"; exit 1; }

# ─── 2) VERIFY IN POSTGRES ───────────────────────────────────────────────────
echo "2) Checking user record in Postgres…"
USER_ID=$(psql -h $DB_HOST -U $DB_USER -d $DB_NAME -Atc \
  "SELECT id FROM users WHERE email = '$EMAIL';")

if [[ -n "$USER_ID" ]]; then
  echo " → Found user id=$USER_ID"
else
  echo " → User not found in database"; exit 1
fi

# ─── 3) VERIFY SEAWEEDFS FILER ───────────────────────────────────────────────
echo "3) Checking SeaweedFS filer health…"
if curl -sf "http://${SEAWEEDFS_FILER}/" >/dev/null; then
  echo " → SeaweedFS filer is reachable"
else
  echo " → SeaweedFS filer is not reachable"; exit 1
fi

echo "✅ All integration checks passed!"
