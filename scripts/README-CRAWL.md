# Hướng dẫn chạy Crawl dữ liệu

## 1. Khởi động Database (Docker)

Từ thư mục gốc project:

```bash
cd docker
docker compose up -d
```

Sẽ chạy:
- **Redis** (port 6379) – queue cho crawl
- **MySQL** (port 3306) – lưu stories, chapter, images

Đợi MySQL sẵn sàng (khoảng 20–30 giây lần đầu):

```bash
docker compose ps
# Hoặc: mysqladmin -h 127.0.0.1 -P 3306 -u root -pcrawl_secret ping
```

## 2. Biến môi trường (tuỳ chọn)

Copy và chỉnh nếu cần:

```bash
cp .env.example .env
# Chỉnh MYSQL_DSN, REDIS_ADDR nếu không dùng localhost
```

Mặc định (không cần .env):
- `MYSQL_DSN=root:crawl_secret@tcp(localhost:3306)/story?parseTime=true`
- `REDIS_ADDR=localhost:6379`

## 3. Chạy Sink Worker (ghi dữ liệu vào MySQL)

Sink worker sẽ **tự chạy migration** (tạo bảng) khi khởi động, rồi consume Redis và ghi vào MySQL.

```bash
cd sink_worker
go run .
```

Giữ terminal này chạy. Worker đọc 3 stream: `series_queue`, `chapter_queue`, `images_queue`.

## 4. Chạy Catalog Worker (crawl và đẩy job vào Redis)

Mở terminal mới:

```bash
cd catalog_worker
go run .
```

Giữ terminal này chạy. Worker đọc: `domain_queue` → crawl domain; `series_queue` → crawl series; `chapter_queue` → crawl chapter/images.

## 5. Seed domain để bắt đầu crawl

Sau khi cả hai worker đang chạy, gửi một domain vào `domain_queue`:

**Cách 1 – Dùng Go (không cần redis-cli):**

```bash
cd catalog_worker
go run ./cmd/seed
# Hoặc domain khác:
go run ./cmd/seed "https://www.webtoons.com/en/"
```

**Cách 2 – Script (cần redis-cli hoặc Docker):**

```bash
./scripts/seed-crawl.sh
# Hoặc: ./scripts/seed-crawl.sh "https://www.webtoons.com/en/"
```

**Cách 3 – redis-cli trực tiếp:**

```bash
redis-cli XADD domain_queue '*' domain_url "https://www.webtoons.com/en/"
# Hoặc Redis trong Docker:
docker exec crawl-redis redis-cli XADD domain_queue '*' domain_url "https://www.webtoons.com/en/"
```

Luồng dữ liệu:
1. **domain_queue** → catalog_worker crawl trang chủ → đẩy vào **series_queue**
2. **series_queue** → sink_worker ghi story vào MySQL; catalog_worker crawl từng series → đẩy vào **chapter_queue**
3. **chapter_queue** → sink_worker ghi chapter; catalog_worker crawl ảnh → đẩy vào **images_queue**
4. **images_queue** → sink_worker ghi images vào MySQL

## 6. Chạy API (xem dữ liệu đã crawl)

```bash
cd server-go
go run ./cmd/server
```

API: http://localhost:8081  
Ví dụ: `GET /api/v1/stories`, `GET /api/v1/stories/:id`, `GET /api/v1/chapter/:storyId/:chapterId` (danh sách ảnh chương).

## Tóm tắt thứ tự chạy

| Bước | Lệnh | Ghi chú |
|------|------|--------|
| 1 | `cd docker && docker compose up -d` | Redis + MySQL, có volume lưu dữ liệu |
| 2 | `cd sink_worker && go run .` | Migration tự chạy, consume 3 queue |
| 3 | `cd catalog_worker && go run .` | Crawl và đẩy job |
| 4 | `./scripts/seed-crawl.sh` | Gửi domain vào domain_queue |
| 5 | `cd server-go && go run ./cmd/server` | (Tuỳ chọn) API xem dữ liệu |
