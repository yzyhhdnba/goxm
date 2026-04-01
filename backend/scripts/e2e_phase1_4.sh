#!/usr/bin/env bash
set -euo pipefail

API='http://127.0.0.1:8080/api/v1'
CREATOR_USER='creator_e2e_0331'
CREATOR_PWD='Creator1234'
ADMIN_USER='hhd'
ADMIN_PWD='Admin1234'
TITLE="E2E上传视频_$(date +%H%M%S)"

creator_token=$(curl -s -X POST "$API/auth/login" -H 'Content-Type: application/json' -d "{\"username\":\"$CREATOR_USER\",\"password\":\"$CREATOR_PWD\"}" | jq -r '.data.access_token // empty')
admin_token=$(curl -s -X POST "$API/auth/login" -H 'Content-Type: application/json' -d "{\"username\":\"$ADMIN_USER\",\"password\":\"$ADMIN_PWD\"}" | jq -r '.data.access_token // empty')

if [[ -z "$creator_token" || -z "$admin_token" ]]; then
  echo "[fatal] token missing"
  exit 1
fi

echo "== phase1 users/me =="
curl -s "$API/users/me" -H "Authorization: Bearer $creator_token" | jq '{code,message,user:.data.username,role:.data.role}'

echo "== phase2 feed/recommend =="
curl -s "$API/feed/recommend?limit=6" | jq '{code,message,count:(.data.list|length)}'

echo "== phase2 video detail =="
curl -s "$API/videos/1" -H "Authorization: Bearer $creator_token" | jq '{code,message,id:.data.id,title:.data.title}'

echo "== phase3 comment create =="
curl -s -X POST "$API/videos/1/comments" -H 'Content-Type: application/json' -H "Authorization: Bearer $creator_token" -d '{"content":"E2E评论-Phase1-4联调验证"}' | jq '{code,message,id:(.data.id // null)}'

create_video_resp=$(curl -s -X POST "$API/videos" -H 'Content-Type: application/json' -H "Authorization: Bearer $creator_token" -d "{\"area_id\":1,\"title\":\"$TITLE\",\"description\":\"E2E upload pipeline test\"}")
video_id=$(echo "$create_video_resp" | jq -r '.data.id // empty')

echo "== phase4 create video =="
echo "$create_video_resp" | jq '{code,message,id:.data.id,title:.data.title,review_status:.data.review_status}'

if [[ -z "$video_id" ]]; then
  echo "[fatal] create video failed"
  exit 1
fi

echo "== phase4 upload source =="
curl -s -X POST "$API/videos/$video_id/source" -H "Authorization: Bearer $creator_token" -F "file=@/Users/hhd/Desktop/test/goxm/frontend/src/assets/video.webm" | jq '{code,message,video_id:.data.video_id,source_path:.data.source_path}'

echo "== phase4 upload cover =="
curl -s -X POST "$API/videos/$video_id/cover" -H "Authorization: Bearer $creator_token" -F "file=@/Users/hhd/Desktop/test/goxm/frontend/src/assets/head.jpg" | jq '{code,message,video_id:.data.video_id,cover_url:.data.cover_url}'

echo "== phase4 creator list =="
curl -s "$API/creator/videos?review_status=all&page=1&page_size=20" -H "Authorization: Bearer $creator_token" | jq '{code,message,count:(.data.list|length)}'

echo "== phase4 admin pending =="
curl -s "$API/admin/videos/pending?page=1&page_size=20" -H "Authorization: Bearer $admin_token" | jq '{code,message,count:(.data.list|length)}'

echo "== phase4 admin approve =="
curl -s -X POST "$API/admin/videos/$video_id/approve" -H "Authorization: Bearer $admin_token" | jq '{code,message,id:.data.id,review_status:.data.review_status}'

kw=$(printf '%s' "$TITLE" | jq -sRr @uri)
echo "== phase4 search videos =="
curl -s "$API/search/videos?keyword=$kw&page=1&page_size=20" | jq '{code,message,count:(.data.list|length),first_title:(.data.list[0].title // null)}'

echo "== phase4 history report =="
curl -s -X POST "$API/histories" -H 'Content-Type: application/json' -H "Authorization: Bearer $creator_token" -d "{\"video_id\":$video_id,\"progress_seconds\":17}" | jq '{code,message}'

echo "== phase4 history list =="
curl -s "$API/histories?page=1&page_size=20" -H "Authorization: Bearer $creator_token" | jq '{code,message,count:(.data.list|length),first_video_id:(.data.list[0].video_id // null)}'
