# Chạy translate-worker trên Windows (GPU RTX 3060 Ti)

## Yêu cầu

- **Windows 10/11** (64-bit)
- **Python 3.10+** — cài từ [python.org](https://www.python.org/downloads/), tick "Add Python to PATH"
- **NVIDIA Driver** — bản mới (tải từ [nvidia.com/drivers](https://www.nvidia.com/Download/index.aspx))
- **CUDA Toolkit** (tùy chọn, thường đã kèm driver mới) — [CUDA 12.x](https://developer.nvidia.com/cuda-downloads) nếu cần

## Bước 1: Cài đặt (chạy một lần)

Mở **Command Prompt** hoặc **PowerShell** trong thư mục `translate-worker`:

```cmd
cd translate-worker
scripts\install-windows-gpu.bat
```

Script sẽ:

1. Tạo virtual environment `venv` (nếu chưa có)
2. Cài **PyTorch 2.6+ với CUDA 12.4** (cu124 — yêu cầu bởi Transformers CVE-2025-32434)
3. Cài các package còn lại (transformers, flask, sacremoses, …)

Nếu lỗi với cu124, script sẽ thử cu121 (PyTorch 2.5). Khi đó có thể gặp lỗi `torch.load` — cần cài torch 2.6+ từ PyPI (sẽ dùng CPU) hoặc dùng cu124 khi có bản wheel.

```cmd
venv\Scripts\activate
pip install "torch>=2.6" torchvision torchaudio --index-url https://download.pytorch.org/whl/cu124
pip install -r requirements.txt
```

## Bước 2: Chạy service

**Cách 1 — File .bat (khuyên dùng):**

```cmd
scripts\run-translate-worker.bat
```

**Cách 2 — PowerShell:**

```powershell
.\scripts\run-translate-worker.ps1
```

**Cách 3 — Tự gõ lệnh:**

```cmd
venv\Scripts\activate
set PYTORCH_CUDA_ALLOC_CONF=expandable_segments:True
python -m translate_worker
```

Service chạy tại **http://localhost:8082**. Log sẽ ghi `Model loaded (device=cuda)` nếu đang dùng GPU.

## Kiểm tra GPU

```cmd
venv\Scripts\python.exe -c "import torch; print('CUDA:', torch.cuda.is_available()); print('GPU:', torch.cuda.get_device_name(0) if torch.cuda.is_available() else 'N/A')"
```

Kết quả mong đợi: `CUDA: True`, `GPU: NVIDIA GeForce RTX 3060 Ti` (hoặc tên tương tự).

## Test nhanh

```cmd
curl -X POST http://localhost:8082/translate -H "Content-Type: application/json" -d "{\"text\": \"Hello world\"}"
```

Hoặc dùng Postman / app iOS gọi `POST http://<IP-may>:8082/translate` với body `{"text": "..."}`.

## Lỗi thường gặp

| Lỗi | Cách xử lý |
|-----|------------|
| `CUDA out of memory` | Giảm batch hoặc đóng app khác dùng GPU; có thể set `CUDA_VISIBLE_DEVICES=0` |
| `torch.cuda.is_available() is False` | Cài lại PyTorch với `cu121` hoặc `cu118`; cập nhật NVIDIA Driver |
| Script PowerShell bị chặn | Chạy: `Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned` |
