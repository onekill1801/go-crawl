"""Dịch văn bản bằng model AI (Hugging Face Transformers) — chạy local trên máy."""

from __future__ import annotations

import logging
from typing import Optional

from .config import DEFAULT_MODEL

logger = logging.getLogger(__name__)

_model = None
_tokenizer = None

# Cache kết quả cho text ngắn (tránh dịch lặp lại cùng câu)
_cache: dict[str, str] = {}
_CACHE_MAX = 500  # số câu tối đa lưu cache


def _get_model_and_tokenizer():
    global _model, _tokenizer
    if _model is None:
        from transformers import AutoModelForSeq2SeqLM, AutoTokenizer
        import torch

        logger.info("Loading translation model: %s", DEFAULT_MODEL)
        _tokenizer = AutoTokenizer.from_pretrained(DEFAULT_MODEL)
        _model = AutoModelForSeq2SeqLM.from_pretrained(DEFAULT_MODEL)
        _model.eval()
        device = "cuda" if torch.cuda.is_available() else "cpu"
        _model = _model.to(device)
        logger.info("Model loaded (device=%s).", device)
    return _model, _tokenizer


def preload_model() -> None:
    """Gọi khi khởi động server để request đầu không phải đợi load model."""
    _get_model_and_tokenizer()
    logger.info("Model preloaded and ready.")


def translate(
    text: str,
    *,
    max_length: Optional[int] = None,
    truncation: bool = True,
) -> str:
    """Dịch một đoạn văn bản. Trả về chuỗi đã dịch."""
    if not text or not text.strip():
        return text

    text = text.strip()
    # Cache cho text ngắn (giảm inference lặp)
    if len(text) < 200 and text in _cache:
        return _cache[text]

    import torch

    model, tokenizer = _get_model_and_tokenizer()
    device = next(model.parameters()).device

    # max_length vừa đủ theo độ dài input (tránh generate thừa, chậm)
    max_in = min(512, len(text) + 50)
    max_out = max_length or min(512, max_in + 20)

    inputs = tokenizer(
        text,
        return_tensors="pt",
        truncation=truncation,
        max_length=max_in,
        padding=True,
    ).to(device)

    with torch.inference_mode():
        out = model.generate(**inputs, max_length=max_out, num_beams=1)

    decoded = tokenizer.decode(out[0], skip_special_tokens=True).strip()

    if len(text) < 200 and len(_cache) < _CACHE_MAX:
        _cache[text] = decoded

    return decoded
