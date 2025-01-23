@echo off
:: 設置 proto 文件目錄
set PROTO_PATH=proto

:: 設置輸出目錄
set OUTPUT_PATH=generated

:: proto 檔案名稱
set NAME=ServiceExample.proto

:: 確保輸出目錄存在
if not exist %OUTPUT_PATH% (
    mkdir %OUTPUT_PATH%
)

:: 生成 Go 代碼，禁用目錄結構
protoc --proto_path=%PROTO_PATH% --go_out=%OUTPUT_PATH% --go_opt=paths=source_relative %PROTO_PATH%\%NAME%

:: 生成 C# 代碼，直接輸出
protoc --proto_path=%PROTO_PATH% --csharp_out=%OUTPUT_PATH% %PROTO_PATH%\%NAME%

:: 暫停以查看輸出
pause
