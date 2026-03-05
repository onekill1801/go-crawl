# translate-worker

Service dịch văn bản bằng model AI (Hugging Face), chạy dưới dạng **HTTP API**. Không dùng Redis.

## Môi trường ảo (venv)

```bash
cd translate-worker
python -m venv .venv
source .venv/bin/activate   # Windows: .venv\Scripts\activate
pip install -r requirements.txt
```

## Chạy service

```bash
python -m translate_worker
```

Mặc định lắng nghe tại **http://0.0.0.0:8082**.

## API

| Method | Path | Mô tả |
|--------|------|--------|
| GET | `/health` | Health check |
| POST | `/translate` | Dịch văn bản |

**POST /translate**

- Body (JSON): `{"text": "Hello world"}` hoặc `{"content": "..."}`
- Trả về: `{"original": "...", "translated": "..."}`

Ví dụ:

```bash
curl -X POST http://localhost:8082/translate \
  -H "Content-Type: application/json" \
  -d '{"text": "Hello, how are you?"}'
```

## Biến môi trường

| Biến | Mặc định | Mô tả |
|------|----------|--------|
| `HOST` | `0.0.0.0` | Bind address |
| `PORT` | `8082` | Cổng HTTP |
| `TRANSLATE_MODEL` | `Helsinki-NLP/opus-mt-en-vi` | Model dịch (en → vi) |

Copy `.env.example` thành `.env` và sửa nếu cần.

## Tốc độ

- **Lần chạy service**: model được load ngay khi start (preload), request đầu không bị chậm vì đợi load.
- **Chạy trên CPU**: mỗi câu thường ~0.5–2 giây tùy độ dài. Dùng **GPU** (CUDA) sẽ nhanh hơn nhiều.
- **Cache**: cùng một đoạn text ngắn (&lt;200 ký tự) dịch lần 2 trở đi trả về ngay từ cache.
