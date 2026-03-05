#!/usr/bin/env bash
# Khởi chạy Docker (Redis + MySQL), sink_worker, catalog_worker và seed domain.
# Dùng cho việc tạo dữ liệu cho app.

set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "==> 1. Docker (Redis + MySQL)..."
cd docker
docker compose up -d
cd "$ROOT"

echo "==> Đợi MySQL sẵn sàng (~25s)..."
sleep 25

echo "==> 2. Sink worker (consume Redis -> insert MySQL) - chạy nền..."
cd sink_worker
go run . &
SINK_PID=$!
cd "$ROOT"

sleep 5
echo "==> 3. Catalog worker (crawl -> đẩy job vào Redis) - chạy nền..."
cd catalog_worker
go run . &
CATALOG_PID=$!
cd "$ROOT"

sleep 4
echo "==> 4. Seed domain_queue..."
cd catalog_worker
go run ./cmd/seed
cd "$ROOT"

echo ""
echo "✅ Các service đã chạy:"
echo "   - Docker: Redis (6379), MySQL (3306)"
echo "   - Sink worker PID: $SINK_PID"
echo "   - Catalog worker PID: $CATALOG_PID"
echo ""
echo "Dừng: kill $SINK_PID $CATALOG_PID"
echo "API (xem dữ liệu): cd server-go && go run ./cmd/server  -> http://localhost:8081"
