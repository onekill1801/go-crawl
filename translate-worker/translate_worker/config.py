"""Cấu hình từ biến môi trường."""

import os

# Load .env nếu có (khi chạy từ thư mục translate-worker)
try:
    from dotenv import load_dotenv
    _root = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    load_dotenv(os.path.join(_root, ".env"))
except ImportError:
    pass

# HTTP server
HOST = os.getenv("HOST", "0.0.0.0")
PORT = int(os.getenv("PORT", "8082"))

# Model chạy local: tự tải lần đầu (from_pretrained). Mặc định VinAI (chất lượng tốt).
# Nếu VinAI lỗi tokenizer sẽ tự fallback sang OPUS-MT.
DEFAULT_MODEL = os.getenv(
    "TRANSLATE_MODEL",
    "vinai/vinai-translate-en2vi-v2",
)
SOURCE_LANG = os.getenv("SOURCE_LANG", "en")
TARGET_LANG = os.getenv("TARGET_LANG", "vi")

# Tùy chọn: gọi HF Inference API thay vì chạy local (đặt USE_HF_INFERENCE_API=1 + HF_TOKEN).
USE_HF_INFERENCE_API = os.getenv("USE_HF_INFERENCE_API", "").lower() in ("1", "true", "yes")
HF_TOKEN = os.getenv("HF_TOKEN", "")
HF_TRANSLATION_MODEL = os.getenv("HF_TRANSLATION_MODEL", "vinai/vinai-translate-en2vi-v2")

# Beam search local: 8 = chất lượng cao, chậm hơn
NUM_BEAMS = int(os.getenv("NUM_BEAMS", "8"))
LENGTH_PENALTY = float(os.getenv("LENGTH_PENALTY", "1.2"))
# 0 = tắt (khuyên dùng). 3–4 có thể giảm lặp nhưng đôi khi làm tệ bản dịch.
NO_REPEAT_NGRAM_SIZE = int(os.getenv("NO_REPEAT_NGRAM_SIZE", "0"))
