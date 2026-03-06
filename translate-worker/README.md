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

## Chạy local — model tự tải

- **Mặc định:** model **VinAI** (`vinai/vinai-translate-en2vi-v2`) chạy **local**. Lần đầu chạy service sẽ **tự tải** model về (khoảng ~1GB), lần sau dùng lại từ cache.
- **Nếu VinAI bị lỗi tokenizer** (một số môi trường), service sẽ **tự chuyển** sang model **OPUS-MT** (`Helsinki-NLP/opus-mt-en-vi`). Có thể đặt sẵn trong `.env`: `TRANSLATE_MODEL=Helsinki-NLP/opus-mt-en-vi`.
- Không cần API key; mọi thứ chạy trên máy. Tùy chọn: bật `USE_HF_INFERENCE_API=1` và `HF_TOKEN` nếu muốn gọi API thay vì chạy local.

## Biến môi trường

| Biến | Mặc định | Mô tả |
|------|----------|--------|
| `HOST` | `0.0.0.0` | Bind address |
| `PORT` | `8082` | Cổng HTTP |
| `TRANSLATE_MODEL` | `vinai/vinai-translate-en2vi-v2` | Model local; lần đầu tự tải |
| `NUM_BEAMS` | `8` | Số beam |
| `LENGTH_PENALTY` | `1.2` | Penalty độ dài câu |
| `NO_REPEAT_NGRAM_SIZE` | `0` | 0 = tắt |

Copy `.env.example` thành `.env` và sửa nếu cần.

## Tốc độ

- **Lần đầu chạy**: model được tải từ Hugging Face (nếu chưa có trong cache), có thể mất vài phút.
- **Các lần sau**: model load từ cache khi start service.
- **Chạy trên CPU**: mỗi câu ~0.5–2 giây. **GPU** (CUDA) nhanh hơn nhiều.
- **Cache**: cùng đoạn text ngắn (&lt;200 ký tự) dịch lần 2 trở đi trả về ngay từ cache.
