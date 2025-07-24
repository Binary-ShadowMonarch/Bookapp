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

# MinIO + mc client (mc must be installed:
# https://docs.min.io/docs/minio-client-quickstart-guide)
MC_ALIAS="local"
MINIO_ENDPOINT="localhost:9000"
MINIO_KEY="admin"
MINIO_SECRET="secretpassword"

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

# ─── 3) VERIFY MINIO BUCKET ──────────────────────────────────────────────────
echo "3) Checking bucket user-${USER_ID} in MinIO…"
mc alias set $MC_ALIAS http://${MINIO_ENDPOINT} $MINIO_KEY $MINIO_SECRET --api S3v4 >/dev/null

if mc ls $MC_ALIAS | grep -q "user-${USER_ID}/"; then
  echo " → Bucket user-${USER_ID} exists"
else
  echo " → Bucket user-${USER_ID} is missing"; exit 1
fi

echo "✅ All integration checks passed!"
