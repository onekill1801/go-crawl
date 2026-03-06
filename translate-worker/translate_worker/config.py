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

# --- NLLB-200-distilled-600M + CTranslate2 (mặc định) ---
# Model CTranslate2: repo Hugging Face hoặc đường dẫn thư mục local. Lần đầu sẽ tự tải.
CT2_MODEL = os.getenv("CT2_MODEL", "JustFrederik/nllb-200-distilled-600M-ct2")
# Tokenizer NLLB (dùng chung với model gốc Facebook)
NLLB_TOKENIZER = os.getenv("NLLB_TOKENIZER", "facebook/nllb-200-distilled-600M")
# Mã ngôn ngữ NLLB: eng_Latn (Anh), vie_Latn (Việt), fra_Latn (Pháp), ...
NLLB_SOURCE_LANG = os.getenv("NLLB_SOURCE_LANG", "eng_Latn")
NLLB_TARGET_LANG = os.getenv("NLLB_TARGET_LANG", "vie_Latn")
# Thiết bị: cpu, cuda, auto
CT2_DEVICE = os.getenv("CT2_DEVICE", "auto")
# Kiểu tính toán: default, float16, int8, int8_float16, ...
CT2_COMPUTE_TYPE = os.getenv("CT2_COMPUTE_TYPE", "default")
# Tham số dịch (ưu trong .env)
CT2_BEAM_SIZE = int(os.getenv("CT2_BEAM_SIZE", "4"))
CT2_LENGTH_PENALTY = float(os.getenv("CT2_LENGTH_PENALTY", "1.0"))
CT2_PATIENCE = float(os.getenv("CT2_PATIENCE", "1.0"))
CT2_MAX_DECODING_LENGTH = int(os.getenv("CT2_MAX_DECODING_LENGTH", "256"))
CT2_MAX_INPUT_LENGTH = int(os.getenv("CT2_MAX_INPUT_LENGTH", "1024"))
CT2_REPETITION_PENALTY = float(os.getenv("CT2_REPETITION_PENALTY", "1.0"))
CT2_NO_REPEAT_NGRAM_SIZE = int(os.getenv("CT2_NO_REPEAT_NGRAM_SIZE", "0"))
# Số luồng CPU (0 = mặc định)
CT2_INTRA_THREADS = int(os.getenv("CT2_INTRA_THREADS", "0"))
CT2_INTER_THREADS = int(os.getenv("CT2_INTER_THREADS", "1"))
