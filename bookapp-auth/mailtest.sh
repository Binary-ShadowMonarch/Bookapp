#!/bin/bash

# Replace with your actual email and server URL
EMAIL="saurabzz3166@gmail.com"
PASSWORD="password123"
SERVER="http://localhost:8080"

echo "=== Testing Email Verification ==="

# Step 1: Request verification code
echo "1. Requesting verification code..."
curl -X POST "$SERVER/signup/request" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "mail=$EMAIL&password=$PASSWORD"

echo -e "\n"

# Step 2: Check your email for the 6-digit code, then run this:
echo "2. Enter the 6-digit code you received:"
read -p "Code: " CODE

echo "Verifying code..."
curl -X POST "$SERVER/signup/verify" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "mail=$EMAIL&code=$CODE"

echo -e "\n"

# Additional test commands:

echo "=== Additional Tests ==="

# Test wrong method
echo "3. Testing wrong method (should fail):"
curl -X GET "$SERVER/signup/request"
echo -e "\n"

# Test short password
echo "4. Testing short password (should fail):"
curl -X POST "$SERVER/signup/request" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "mail=$EMAIL&password=short"
echo -e "\n"

# Test wrong code
echo "5. Testing wrong code (should fail):"
curl -X POST "$SERVER/signup/verify" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "mail=$EMAIL&code=000000"
echo -e "\n"