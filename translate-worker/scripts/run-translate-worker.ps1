# Translate-Worker - PowerShell (Windows, GPU)
# Chay: .\scripts\run-translate-worker.ps1
# Hoac: powershell -ExecutionPolicy Bypass -File .\scripts\run-translate-worker.ps1

$ErrorActionPreference = "Stop"
$ProjectRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
Set-Location $ProjectRoot

$python = Join-Path $ProjectRoot ".venv", "Scripts", "python.exe"
if (-not (Test-Path $python)) {
    Write-Host "LOI: Chua co .venv. Chay truoc: scripts\install-windows-gpu.bat" -ForegroundColor Red
    exit 1
}

$env:PYTORCH_CUDA_ALLOC_CONF = "expandable_segments:True"
Write-Host "Translate-Worker (GPU) - http://localhost:8082" -ForegroundColor Cyan
Write-Host "Thoat: Ctrl+C" -ForegroundColor Gray
& $python -m translate_worker
