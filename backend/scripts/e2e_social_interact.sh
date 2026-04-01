#!/usr/bin/env bash
set -euo pipefail

API='http://127.0.0.1:8080/api/v1'
U='social_e2e_0331'
E='social_e2e_0331@example.com'
P='Social1234'

curl -s -X POST "$API/auth/register" \
  -H 'Content-Type: application/json' \
  -d "{\"username\":\"$U\",\"email\":\"$E\",\"password\":\"$P\"}" >/dev/null || true

login=$(curl -s -X POST "$API/auth/login" \
  -H 'Content-Type: application/json' \
  -d "{\"username\":\"$U\",\"password\":\"$P\"}")

token=$(echo "$login" | jq -r '.data.access_token // empty')
if [[ -z "$token" ]]; then
  echo "[fatal] social token missing"
  echo "$login"
  exit 1
fi

echo "== like video =="
curl -s -X POST "$API/videos/1/likes" -H "Authorization: Bearer $token" | jq '{code,message}'

echo "== like status =="
curl -s "$API/videos/1/likes/me" -H "Authorization: Bearer $token" | jq '{code,message,liked:.data.liked}'

echo "== unlike video =="
curl -s -X DELETE "$API/videos/1/likes" -H "Authorization: Bearer $token" | jq '{code,message}'

echo "== favorite video =="
curl -s -X POST "$API/videos/1/favorites" -H "Authorization: Bearer $token" | jq '{code,message}'

echo "== favorite status =="
curl -s "$API/videos/1/favorites/me" -H "Authorization: Bearer $token" | jq '{code,message,favorited:.data.favorited}'

echo "== unfavorite video =="
curl -s -X DELETE "$API/videos/1/favorites" -H "Authorization: Bearer $token" | jq '{code,message}'

echo "== follow user 3 =="
curl -s -X POST "$API/users/3/follow" -H "Authorization: Bearer $token" | jq '{code,message}'

echo "== follow status user 3 =="
curl -s "$API/users/3/follow-status" -H "Authorization: Bearer $token" | jq '{code,message,followed:.data.followed}'

echo "== followers of user 3 =="
curl -s "$API/users/3/followers?page=1&page_size=20" | jq '{code,message,count:(.data.list|length)}'

echo "== my following list =="
uid=$(curl -s "$API/users/me" -H "Authorization: Bearer $token" | jq -r '.data.id')
curl -s "$API/users/$uid/following?page=1&page_size=20" | jq '{code,message,count:(.data.list|length)}'

echo "== unfollow user 3 =="
curl -s -X DELETE "$API/users/3/follow" -H "Authorization: Bearer $token" | jq '{code,message}'
