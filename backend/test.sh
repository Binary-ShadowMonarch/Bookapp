#!/usr/bin/env bash
set -euo pipefail

API="http://localhost:8080"

echo "1) Register"
curl -s -X POST \
  -d "mail=alice@example.com&password=secret123" \
  "$API/register" && echo -e "\n"

echo "2) Login"
TOKEN=$(curl -s -X POST \
  -d "mail=alice@example.com&password=secret123" \
  "$API/login" | jq -r .token)
echo "  → token: $TOKEN"

echo "3) Protected (valid token)"
curl -s -H "Authorization: Bearer $TOKEN" \
  "$API/protected" && echo -e "\n"

echo "4) Protected (invalid token)"
curl -s -H "Authorization: Bearer WRONG" \
  "$API/protected" && echo -e "\n"

echo "5) Logout (no-op, but should return 204)"
curl -i -X POST \
  -H "Authorization: Bearer $TOKEN" \
  "$API/logout" && echo -e "\n"
