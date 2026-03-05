"""HTTP service: nhận văn bản, dịch bằng model AI, trả JSON."""

from __future__ import annotations

import logging
from typing import Any

from flask import Flask, jsonify, request

from .config import HOST, PORT
from .translator import preload_model, translate

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
)
logger = logging.getLogger(__name__)

app = Flask(__name__)


@app.route("/health", methods=["GET"])
def health() -> tuple[dict[str, str], int]:
    return {"status": "ok"}, 200


@app.route("/translate", methods=["POST"])
def do_translate() -> tuple[dict[str, Any], int]:
    """Body JSON: {"text": "..."} hoặc {"content": "..."}. Trả về {"translated": "..."}."""
    if not request.is_json:
        return jsonify({"error": "Content-Type must be application/json"}), 400

    data = request.get_json() or {}
    text = data.get("text") or data.get("content") or ""
    if isinstance(text, bytes):
        text = text.decode("utf-8", errors="replace")

    if not text.strip():
        return jsonify({"translated": "", "original": ""}), 200

    try:
        translated = translate(text)
        return jsonify({"original": text, "translated": translated}), 200
    except Exception as e:
        logger.exception("Translate error: %s", e)
        return jsonify({"error": str(e), "translated": text}), 500


def run_server() -> None:
    logger.info("Preloading model (request đầu sẽ không bị chậm)...")
    preload_model()
    logger.info("Translate service http://%s:%s", HOST, PORT)
    app.run(host=HOST, port=PORT, threaded=True)
