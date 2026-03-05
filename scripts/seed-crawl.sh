#!/usr/bin/env bash
# Seed domain_queue để bắt đầu crawl (catalog_worker sẽ đọc và crawl danh sách series)
# Cách dùng: ./scripts/seed-crawl.sh [domain_url]
# Mặc định: https://www.webtoons.com/en/

set -e
REDIS_CLI="${REDIS_CLI:-redis-cli}"
REDIS_ADDR="${REDIS_ADDR:-localhost:6379}"
DOMAIN_URL="${1:-https://www.webtoons.com/en/}"

echo "→ Redis: $REDIS_ADDR"
echo "→ Seed domain_queue với domain_url=$DOMAIN_URL"

# Nếu chạy Redis trong Docker
if command -v docker &>/dev/null && docker ps --format '{{.Names}}' 2>/dev/null | grep -q crawl-redis; then
  docker exec crawl-redis redis-cli XADD domain_queue '*' domain_url "$DOMAIN_URL"
  echo "✅ Đã gửi 1 message vào domain_queue (qua Docker)."
else
  # redis-cli: -h host -p port hoặc -u redis://host:port
  if [[ "$REDIS_ADDR" == *:* ]]; then
    host="${REDIS_ADDR%%:*}"
    port="${REDIS_ADDR##*:}"
    "$REDIS_CLI" -h "$host" -p "$port" XADD domain_queue '*' domain_url "$DOMAIN_URL"
  else
    "$REDIS_CLI" -h "$REDIS_ADDR" -p 6379 XADD domain_queue '*' domain_url "$DOMAIN_URL"
  fi
  echo "✅ Đã gửi 1 message vào domain_queue."
fi
