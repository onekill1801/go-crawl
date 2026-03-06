# translate-worker

Service dịch văn bản (en → vi) bằng **NLLB-200-distilled-600M + CTranslate2**, chạy **local** dưới dạng HTTP API.

## Môi trường ảo (venv)

```bash
cd translate-worker
python -m venv venv
# Windows: venv\Scripts\activate
# Linux/macOS: source venv/bin/activate
pip install -r requirements.txt
```

## Chạy service

```bash
python -m translate_worker
```

Mặc định lắng nghe tại **http://0.0.0.0:8082**. Lần đầu chạy sẽ **tự tải** model CTranslate2 và tokenizer NLLB từ Hugging Face.

## API

| Method | Path | Mô tả |
|--------|------|--------|
| GET | `/health` | Health check |
| GET | `/info` | Thông tin model, đường dẫn lưu, ngôn ngữ nguồn/đích |
| POST | `/translate` | Dịch văn bản |

**POST /translate**

- Body (JSON): `{"text": "Hello world"}` hoặc `{"content": "..."}`
- Trả về: `{"original": "...", "translated": "..."}`

## Model: NLLB-200-distilled-600M + CTranslate2

- **Model CTranslate2:** `JustFrederik/nllb-200-distilled-600M-ct2` (repo HF) hoặc đường dẫn thư mục local.
- **Tokenizer:** `facebook/nllb-200-distilled-600M`.
- **Ngôn ngữ:** mặc định `eng_Latn` → `vie_Latn`. Có thể đổi trong `.env` (NLLB hỗ trợ 200 ngôn ngữ).

## Biến môi trường (.env) — ưu tham số tại đây

| Biến | Mặc định | Mô tả |
|------|----------|--------|
| `HOST` | `0.0.0.0` | Bind address |
| `PORT` | `8082` | Cổng HTTP |
| `CT2_MODEL` | `JustFrederik/nllb-200-distilled-600M-ct2` | Model CT2 (repo HF hoặc đường dẫn thư mục) |
| `NLLB_TOKENIZER` | `facebook/nllb-200-distilled-600M` | Tokenizer NLLB |
| `NLLB_SOURCE_LANG` | `eng_Latn` | Mã ngôn ngữ nguồn |
| `NLLB_TARGET_LANG` | `vie_Latn` | Mã ngôn ngữ đích |
| `CT2_DEVICE` | `auto` | `cpu`, `cuda`, `auto` |
| `CT2_COMPUTE_TYPE` | `default` | `default`, `float16`, `int8`, `int8_float16` |
| `CT2_BEAM_SIZE` | `4` | Beam size (1 = greedy) |
| `CT2_LENGTH_PENALTY` | `1.0` | Penalty độ dài |
| `CT2_PATIENCE` | `1.0` | Patience beam search |
| `CT2_MAX_DECODING_LENGTH` | `256` | Độ dài tối đa bản dịch |
| `CT2_MAX_INPUT_LENGTH` | `1024` | Độ dài tối đa input |
| `CT2_REPETITION_PENALTY` | `1.0` | Penalty lặp từ |
| `CT2_NO_REPEAT_NGRAM_SIZE` | `0` | Cấm lặp n-gram (0 = tắt) |
| `CT2_INTRA_THREADS` | `0` | Số luồng OpenMP (0 = mặc định) |
| `CT2_INTER_THREADS` | `1` | Số batch song song |

Copy `.env.example` thành `.env` và chỉnh theo ý (ưu tham số trong `.env`).

## Vị trí lưu model

- Model tải từ Hugging Face nằm trong cache mặc định, ví dụ Windows:  
  `C:\Users\<User>\.cache\huggingface\hub\`  
  Thư mục model dạng: `models--JustFrederik--nllb-200-distilled-600M-ct2`.
- Gọi **GET /info** để xem `model`, `model_path`, `tokenizer`, `cache_dir`.

## Tốc độ

- Lần đầu: tải model + tokenizer từ HF (có thể vài phút).
- CTranslate2 tối ưu inference (CPU/GPU). Có thể dùng bản int8: `JustFrederik/nllb-200-distilled-600M-ct2-int8` (đặt `CT2_MODEL` trong `.env`).
