"""Dịch văn bản bằng NLLB-200-distilled-600M + CTranslate2 — chạy local."""

from __future__ import annotations

import logging
import os
from typing import Optional

from .config import (
    CT2_BEAM_SIZE,
    CT2_COMPUTE_TYPE,
    CT2_DEVICE,
    CT2_INTER_THREADS,
    CT2_INTRA_THREADS,
    CT2_LENGTH_PENALTY,
    CT2_MAX_DECODING_LENGTH,
    CT2_MAX_INPUT_LENGTH,
    CT2_MODEL,
    CT2_NO_REPEAT_NGRAM_SIZE,
    CT2_PATIENCE,
    CT2_REPETITION_PENALTY,
    NLLB_SOURCE_LANG,
    NLLB_TARGET_LANG,
    NLLB_TOKENIZER,
)

logger = logging.getLogger(__name__)

_translator = None
_tokenizer = None
_loaded_model_path: Optional[str] = None

# Cache kết quả cho text ngắn
_cache: dict[str, str] = {}
_CACHE_MAX = 500


def _resolve_ct2_model_path(model_spec: str) -> str:
    """Trả về đường dẫn thư mục model: nếu là repo HF thì tải về cache."""
    if os.path.isdir(model_spec):
        return os.path.abspath(model_spec)
    if "/" in model_spec and "\\" not in model_spec:
        try:
            from huggingface_hub import snapshot_download
            path = snapshot_download(repo_id=model_spec, local_files_only=False)
            return path
        except Exception as e:
            logger.warning("snapshot_download failed for %s: %s", model_spec, e)
    return model_spec


def _get_translator_and_tokenizer():
    global _translator, _tokenizer, _loaded_model_path
    if _translator is None:
        import ctranslate2
        from transformers import AutoTokenizer

        model_path = _resolve_ct2_model_path(CT2_MODEL)
        _loaded_model_path = model_path
        if CT2_DEVICE == "auto":
            try:
                import torch
                device = "cuda" if torch.cuda.is_available() else "cpu"
            except Exception:
                device = "cpu"
        else:
            device = CT2_DEVICE
        compute_type = (CT2_COMPUTE_TYPE or "default").strip() or "default"
        logger.info("Loading NLLB + CTranslate2 (first run may download): %s", CT2_MODEL)
        _translator = ctranslate2.Translator(
            model_path,
            device=device,
            compute_type=compute_type,
            inter_threads=CT2_INTER_THREADS,
            intra_threads=CT2_INTRA_THREADS or 0,
        )
        _tokenizer = AutoTokenizer.from_pretrained(
            NLLB_TOKENIZER,
            src_lang=NLLB_SOURCE_LANG,
            clean_up_tokenization_spaces=True,
        )
        logger.info("Model loaded (device=%s). %s -> %s", device, NLLB_SOURCE_LANG, NLLB_TARGET_LANG)
    return _translator, _tokenizer


def get_model_info() -> dict[str, str]:
    """Trả về thông tin model đang dùng và thư mục lưu."""
    try:
        from huggingface_hub.constants import HF_HUB_CACHE
        cache_dir = os.path.expanduser(os.path.expandvars(HF_HUB_CACHE))
    except Exception:
        cache_dir = os.path.expanduser(
            os.environ.get("HF_HUB_CACHE", os.path.join(os.path.expanduser("~"), ".cache", "huggingface", "hub"))
        )
    model_name = _loaded_model_path or CT2_MODEL
    return {
        "model": CT2_MODEL,
        "model_path": model_name,
        "tokenizer": NLLB_TOKENIZER,
        "source_lang": NLLB_SOURCE_LANG,
        "target_lang": NLLB_TARGET_LANG,
        "cache_dir": cache_dir,
    }


def preload_model() -> None:
    """Tải model local (tự download lần đầu), sẵn sàng dịch."""
    _get_translator_and_tokenizer()
    info = get_model_info()
    logger.info(
        "Model preloaded and ready. Model: %s | Path: %s",
        info["model"],
        info["model_path"],
    )


def translate(
    text: str,
    *,
    max_length: Optional[int] = None,
    truncation: bool = True,
) -> str:
    """Dịch một đoạn văn bản (en -> vi mặc định). Trả về chuỗi đã dịch."""
    if not text or not text.strip():
        return text

    text = text.strip()
    if len(text) < 200 and text in _cache:
        return _cache[text]

    translator, tokenizer = _get_translator_and_tokenizer()
    max_dec = max_length or CT2_MAX_DECODING_LENGTH

    input_ids = tokenizer.encode(text, truncation=truncation, max_length=CT2_MAX_INPUT_LENGTH)
    source_tokens = tokenizer.convert_ids_to_tokens(input_ids)
    target_prefix = [[NLLB_TARGET_LANG]]

    results = translator.translate_batch(
        [source_tokens],
        target_prefix=target_prefix,
        beam_size=CT2_BEAM_SIZE,
        length_penalty=CT2_LENGTH_PENALTY,
        patience=CT2_PATIENCE,
        max_decoding_length=max_dec,
        max_input_length=CT2_MAX_INPUT_LENGTH,
        repetition_penalty=CT2_REPETITION_PENALTY,
        no_repeat_ngram_size=CT2_NO_REPEAT_NGRAM_SIZE,
    )

    # hypotheses[0] là list token; bỏ token đầu (mã ngôn ngữ đích)
    out_tokens = results[0].hypotheses[0][1:]
    decoded = tokenizer.decode(tokenizer.convert_tokens_to_ids(out_tokens), skip_special_tokens=True).strip()

    if len(text) < 200 and len(_cache) < _CACHE_MAX:
        _cache[text] = decoded

    return decoded
