@echo off
REM OpenAPI 文档验证脚本 (Windows版本)

echo 正在验证 OpenAPI 文档...

REM 检查文件是否存在
if not exist "openapi.json" (
    echo ❌ 错误: openapi.json 文件不存在
    exit /b 1
)

REM 基本的 JSON 格式验证
echo 🔍 验证 JSON 格式...
powershell -Command "Get-Content openapi.json | ConvertFrom-Json" >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ JSON 格式正确
) else (
    echo ❌ JSON 格式错误
    exit /b 1
)

REM 检查必需字段
echo 🔍 检查必需字段...

REM 检查 openapi 字段
powershell -Command "try { $json = Get-Content openapi.json | ConvertFrom-Json; if ($json.openapi) { exit 0 } else { exit 1 } } catch { exit 1 }" >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ 包含必需字段: openapi
) else (
    echo ❌ 缺少必需字段: openapi
    exit /b 1
)

REM 检查 info 字段
powershell -Command "try { $json = Get-Content openapi.json | ConvertFrom-Json; if ($json.info) { exit 0 } else { exit 1 } } catch { exit 1 }" >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ 包含必需字段: info
) else (
    echo ❌ 缺少必需字段: info
    exit /b 1
)

REM 检查 paths 字段
powershell -Command "try { $json = Get-Content openapi.json | ConvertFrom-Json; if ($json.paths) { exit 0 } else { exit 1 } } catch { exit 1 }" >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ 包含必需字段: paths
) else (
    echo ❌ 缺少必需字段: paths
    exit /b 1
)

REM 检查 components 字段
powershell -Command "try { $json = Get-Content openapi.json | ConvertFrom-Json; if ($json.components) { exit 0 } else { exit 1 } } catch { exit 1 }" >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ 包含必需字段: components
) else (
    echo ❌ 缺少必需字段: components
    exit /b 1
)

echo 🎉 OpenAPI 文档验证完成！
echo.
echo 文档信息:
powershell -Command "
try {
    $json = Get-Content openapi.json | ConvertFrom-Json
    Write-Host '标题: ' $json.info.title
    Write-Host '版本: ' $json.info.version
    Write-Host '路径数量: ' $json.paths.Count
    Write-Host '组件模式数量: ' $json.components.schemas.Count
} catch {
    Write-Host '无法读取文档信息'
}
"

pause