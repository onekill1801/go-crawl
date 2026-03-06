"""Dịch văn bản bằng model AI (Hugging Face Transformers) — chạy local trên máy."""

from __future__ import annotations

import logging
import os
from typing import Any, Optional

from .config import (
    DEFAULT_MODEL,
    HF_TOKEN,
    HF_TRANSLATION_MODEL,
    LENGTH_PENALTY,
    NO_REPEAT_NGRAM_SIZE,
    NUM_BEAMS,
    USE_HF_INFERENCE_API,
)

logger = logging.getLogger(__name__)

_model = None
_tokenizer = None
# Model mBART (VinAI) cần forced_bos_token_id = token tiếng Việt
_forced_bos_token_id: Optional[int] = None
# Tên model đã load (sau khi có fallback) — dùng cho get_model_info()
_loaded_model_name: Optional[str] = None

# Cache kết quả cho text ngắn (tránh dịch lặp lại cùng câu)
_cache: dict[str, str] = {}
_CACHE_MAX = 500  # số câu tối đa lưu cache


def _get_forced_bos_token_id(tokenizer: Any) -> Optional[int]:
    """Lấy token id ngôn ngữ đích (vi_VN) cho mBART — bắt buộc để dịch đúng."""
    for token in ("vi_VN", "__vi_VN__", "vi"):
        tid = tokenizer.convert_tokens_to_ids(token)
        if tid is not None and tid != tokenizer.unk_token_id and tid != 0:
            return tid
    if getattr(tokenizer, "lang_code_to_id", None):
        return tokenizer.lang_code_to_id.get("vi_VN") or tokenizer.lang_code_to_id.get("vi")
    # mBART-50: Vietnamese thường là 250021
    return 250021


# Model fallback khi VinAI lỗi tokenizer (transformers đọc nhầm SentencePiece bằng tiktoken).
_FALLBACK_MODEL = "Helsinki-NLP/opus-mt-en-vi"


def _get_model_and_tokenizer():
    global _model, _tokenizer, _forced_bos_token_id, _loaded_model_name
    if _model is None:
        from transformers import AutoModelForSeq2SeqLM, AutoTokenizer
        import torch

        model_to_load = DEFAULT_MODEL
        # VinAI: bắt buộc use_fast=False để dùng slow tokenizer (SentencePiece), tránh lỗi convert.
        is_vinai = "vinai" in model_to_load.lower() and "en2vi" in model_to_load.lower()
        for attempt in range(2):
            try:
                if attempt > 0:
                    model_to_load = _FALLBACK_MODEL
                    logger.warning("Fallback to local model: %s", model_to_load)
                logger.info("Loading translation model (first run may download): %s", model_to_load)
                if is_vinai or "vinai" in model_to_load.lower():
                    _prev = os.environ.get("TRANSFORMERS_USE_FAST_TOKENIZER")
                    os.environ["TRANSFORMERS_USE_FAST_TOKENIZER"] = "0"
                    try:
                        _tokenizer = AutoTokenizer.from_pretrained(model_to_load, use_fast=False)
                    finally:
                        if _prev is None:
                            os.environ.pop("TRANSFORMERS_USE_FAST_TOKENIZER", None)
                        else:
                            os.environ["TRANSFORMERS_USE_FAST_TOKENIZER"] = _prev
                else:
                    _tokenizer = AutoTokenizer.from_pretrained(model_to_load)
                _model = AutoModelForSeq2SeqLM.from_pretrained(model_to_load)
                _loaded_model_name = model_to_load
                break
            except Exception as e:
                logger.warning("Load failed for %s: %s", model_to_load, e)
                _model, _tokenizer = None, None
                if attempt == 0 and (is_vinai or "vinai" in (model_to_load or "").lower()):
                    is_vinai = False
                    continue
                raise
        _model.eval()
        device = "cuda" if torch.cuda.is_available() else "cpu"
        _model = _model.to(device)
        if getattr(_model.config, "model_type", None) == "mbart":
            _forced_bos_token_id = _get_forced_bos_token_id(_tokenizer)
            logger.info("mBART: using forced_bos_token_id=%s for Vietnamese.", _forced_bos_token_id)
        else:
            _forced_bos_token_id = None
        logger.info("Model loaded (device=%s).", device)
    return _model, _tokenizer


def get_model_info() -> dict[str, str]:
    """Trả về thông tin model đang dùng và thư mục cache trên máy."""
    try:
        from huggingface_hub import constants
        cache_dir = os.path.expanduser(
            os.environ.get("HF_HUB_CACHE", getattr(constants, "HF_HUB_CACHE", ""))
        )
        if not cache_dir and hasattr(constants, "HF_HUB_CACHE"):
            cache_dir = os.path.expanduser(constants.HF_HUB_CACHE)
    except Exception:
        cache_dir = os.path.expanduser(
            os.environ.get("HF_HUB_CACHE", os.path.join(os.path.expanduser("~"), ".cache", "huggingface", "hub"))
        )
    name = _loaded_model_name or DEFAULT_MODEL
    return {
        "model": name,
        "cache_dir": cache_dir,
        "model_cache_path": os.path.join(cache_dir, "models--" + name.replace("/", "--")),
    }


def preload_model() -> None:
    """Tải model local (tự download lần đầu), rồi sẵn sàng dịch."""
    _get_model_and_tokenizer()
    info = get_model_info()
    logger.info(
        "Model preloaded and ready. Model: %s | Cache: %s",
        info["model"],
        info["cache_dir"],
    )
    logger.info("Model preloaded and ready.")


def _translate_via_hf_api(text: str) -> Optional[str]:
    """Dịch qua Hugging Face Inference API (VinAI trên server) — chất lượng tốt nhất."""
    if not HF_TOKEN or not USE_HF_INFERENCE_API:
        return None
    try:
        import urllib.request
        import json
        url = f"https://api-inference.huggingface.co/models/{HF_TRANSLATION_MODEL}"
        req = urllib.request.Request(
            url,
            data=json.dumps({"inputs": text}).encode("utf-8"),
            headers={
                "Authorization": f"Bearer {HF_TOKEN}",
                "Content-Type": "application/json",
            },
            method="POST",
        )
        with urllib.request.urlopen(req, timeout=60) as resp:
            out = json.loads(resp.read().decode("utf-8"))
        if isinstance(out, str):
            return out.strip() or None
        if isinstance(out, list) and len(out) and isinstance(out[0], dict):
            return (out[0].get("translation_text") or out[0].get("translated_text") or "").strip() or None
        if isinstance(out, dict):
            return (out.get("translation_text") or out.get("translated_text") or "").strip() or None
        return None
    except Exception as e:
        logger.warning("HF Inference API failed: %s", e)
        return None


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

    # Ưu tiên Hugging Face Inference API (VinAI) — chất lượng tốt nhất
    if USE_HF_INFERENCE_API and HF_TOKEN:
        hf_out = _translate_via_hf_api(text)
        if hf_out is not None and hf_out.strip():
            if len(text) < 200 and len(_cache) < _CACHE_MAX:
                _cache[text] = hf_out
            return hf_out
        # API lỗi hoặc trả rỗng -> fallback local
        logger.info("HF API unavailable or empty, using local model.")

    import torch

    model, tokenizer = _get_model_and_tokenizer()
    device = next(model.parameters()).device

    # mBART (VinAI) hỗ trợ max 1024; Marian/OPUS thường 512
    max_src = 1024 if getattr(model.config, "model_type", None) == "mbart" else 512
    max_in = min(max_src, len(text) + 80)
    max_out = max_length or min(max_src, max_in + 60)

    inputs = tokenizer(
        text,
        return_tensors="pt",
        truncation=truncation,
        max_length=max_in,
        padding=True,
    ).to(device)

    generate_kwargs: dict[str, Any] = {
        "max_length": max_out,
        "num_beams": NUM_BEAMS,
        "length_penalty": LENGTH_PENALTY,
        "early_stopping": True,
    }
    if NO_REPEAT_NGRAM_SIZE > 0:
        generate_kwargs["no_repeat_ngram_size"] = NO_REPEAT_NGRAM_SIZE
    if _forced_bos_token_id is not None:
        generate_kwargs["forced_bos_token_id"] = _forced_bos_token_id

    with torch.inference_mode():
        out = model.generate(**inputs, **generate_kwargs)

    decoded = tokenizer.decode(out[0], skip_special_tokens=True).strip()

    if len(text) < 200 and len(_cache) < _CACHE_MAX:
        _cache[text] = decoded

    return decoded
