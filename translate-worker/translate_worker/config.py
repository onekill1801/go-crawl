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

# Model mặc định: English -> Vietnamese (Helsinki-NLP)
DEFAULT_MODEL = os.getenv(
    "TRANSLATE_MODEL",
    "Helsinki-NLP/opus-mt-en-vi",
)
SOURCE_LANG = os.getenv("SOURCE_LANG", "en")
TARGET_LANG = os.getenv("TARGET_LANG", "vi")
