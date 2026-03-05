@echo off
chcp 65001 >nul
setlocal
cd /d "%~dp0.."

echo ========================================
echo  Translate-Worker (Windows - GPU RTX 3060 Ti)
echo ========================================
echo.

if not exist ".venv\Scripts\python.exe" (
    echo LOI: Chua co .venv. Chay truoc: scripts\install-windows-gpu.bat
    pause
    exit /b 1
)

:: Tu dong dung GPU neu PyTorch nhan CUDA
set PYTORCH_CUDA_ALLOC_CONF=expandable_segments:True

echo Khoi dong service (model se load vao GPU khi start)...
echo Service: http://localhost:8082
echo Thoat: Ctrl+C
echo.

.venv\Scripts\python.exe -m translate_worker

pause
