#!/usr/bin/env bash
set -euo pipefail

HOST="http://localhost:8080"
EMAIL="poudelsaurab04@gmail.com"   # your test user
PASSWORD="Saurab!@#"
COOKIEJAR="./cookies.txt"
PROTECTED_PATH="/protected/library"
ACCESS_TTL=25  # seconds

# do_request <curl args…> → prints "<code>:<body>"
do_request(){
  local resp body code
  resp=$(curl -s -w "\n%{http_code}" "$@")
  body=$(printf '%s\n' "$resp" | sed '$d')
  code=$(printf '%s\n' "$resp" | tail -n1)
  printf "%s:%s\n" "$code" "$body"
}

echo
echo "1) Login **without** creds"
do_request \
  -c /dev/null \
  -X POST "$HOST/login" \
  -H "Content-Type: application/x-www-form-urlencoded"

echo
echo "2) Login **with** valid creds"
do_request \
  -c "$COOKIEJAR" \
  -X POST "$HOST/login" \
    --data "mail=$EMAIL" \
    --data "password=$PASSWORD"

echo
echo "3) Access protected immediately"
do_request \
  -b "$COOKIEJAR" \
  -H "Accept: application/json" \
  "$HOST$PROTECTED_PATH"

echo
echo "4) Sleeping $((ACCESS_TTL+1))s for access_token expiry…"
sleep $((ACCESS_TTL + 1))

echo
echo "5) Access protected **after** expiry"
do_request \
  -b "$COOKIEJAR" \
  -H "Accept: application/json" \
  "$HOST$PROTECTED_PATH"

echo
echo "6) Refresh tokens"
do_request \
  -b "$COOKIEJAR" \
  -c "$COOKIEJAR" \
  -X POST "$HOST/refresh" \
  -H "Accept: application/json"

echo
echo "7) Access protected again"
do_request \
  -b "$COOKIEJAR" \
  -H "Accept: application/json" \
  "$HOST$PROTECTED_PATH"

echo
echo "✅ Done — all steps executed."
