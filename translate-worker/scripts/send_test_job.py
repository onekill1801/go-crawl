#!/usr/bin/env python3
"""Gửi request dịch thử tới HTTP service translate-worker."""

import json
import os
import sys

# Thêm thư mục cha
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

def main():
    try:
        from urllib.request import urlopen, Request
    except ImportError:
        import urllib.request as u
        urlopen, Request = u.urlopen, u.Request

    from translate_worker.config import HOST, PORT

    url = f"http://{HOST}:{PORT}/translate"
    data = json.dumps({"text": "Hello world, this is a test for translation."}).encode("utf-8")
    req = Request(url, data=data, method="POST", headers={"Content-Type": "application/json"})
    with urlopen(req, timeout=30) as resp:
        out = json.loads(resp.read().decode())
    print("Response:", json.dumps(out, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
