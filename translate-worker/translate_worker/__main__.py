"""Chạy service: python -m translate_worker (HTTP server, không dùng Redis)."""

from .worker import run_server

if __name__ == "__main__":
    run_server()
