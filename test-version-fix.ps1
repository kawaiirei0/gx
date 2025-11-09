# 测试版本检测修复

Write-Host "`n=== Testing Version Detection Fix ===" -ForegroundColor Cyan
Write-Host ""

# 1. 测试 list 命令
Write-Host "1. Testing 'gx list' command..." -ForegroundColor Yellow
.\build\gx.exe list
Write-Host ""

# 2. 测试 current 命令
Write-Host "2. Testing 'gx current' command..." -ForegroundColor Yellow
.\build\gx.exe current
Write-Host ""

# 3. 测试 list --verbose
Write-Host "3. Testing 'gx list --verbose' command..." -ForegroundColor Yellow
.\build\gx.exe list --verbose
Write-Host ""

# 4. 显示配置文件内容
Write-Host "4. Checking config file..." -ForegroundColor Yellow
$configPath = "$env:USERPROFILE\.gx\config.json"
if (Test-Path $configPath) {
    Write-Host "Config file location: $configPath" -ForegroundColor Green
    Get-Content $configPath | ConvertFrom-Json | ConvertTo-Json -Depth 10
} else {
    Write-Host "Config file not found at: $configPath" -ForegroundColor Red
}
Write-Host ""

Write-Host "=== Test Complete ===" -ForegroundColor Cyan
