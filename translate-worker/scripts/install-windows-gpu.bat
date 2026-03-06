@echo off
chcp 65001 >nul
setlocal
cd /d "%~dp0.."

echo ========================================
echo  Translate-Worker - Cai dat (Windows GPU)
echo  GPU: NVIDIA RTX 3060 Ti (CUDA)
echo ========================================
echo.

:: Yeu cau: Python 3.10+ da cai, NVIDIA Driver + CUDA (neu chua co)
:: Tai CUDA Toolkit: https://developer.nvidia.com/cuda-downloads

if not exist "venv" (
    echo [1/3] Tao virtual environment...
    py -m venv venv
    if errorlevel 1 (
        echo LOI: Khong tao duoc venv. Kiem tra da cai Python chua.
        pause
        exit /b 1
    )
) else (
    echo [1/3] Thu muc venv da ton tai.
)

call venv\Scripts\activate.bat

echo.
echo [2/3] Cai PyTorch 2.6+ voi CUDA (yeu cau boi Transformers CVE-2025-32434)...
echo       Thu cu124 truoc (co torch 2.6), neu loi thi thu cu121 (torch 2.5)...
pip install "torch>=2.6" torchvision torchaudio --index-url https://download.pytorch.org/whl/cu124
if errorlevel 1 (
    echo Thu lai voi cu121 (torch 2.5 - can transformers moi hon hoac chay CPU)...
    pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu121
)

echo.
echo [3/3] Cai cac package con lai (khong ghi de torch CUDA)...
pip install flask transformers sentencepiece accelerate python-dotenv sacremoses tiktoken
if errorlevel 1 (
    echo LOI: pip install that bai.
    pause
    exit /b 1
)

echo.
echo Kiem tra GPU...
py -c "import torch; print('CUDA available:', torch.cuda.is_available()); print('Device:', torch.cuda.get_device_name(0) if torch.cuda.is_available() else 'CPU')"
echo.
echo ========================================
echo  Cai dat xong. Chay: scripts\run-translate-worker.bat
echo ========================================
pause
