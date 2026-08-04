@echo off
echo === 编译 ArkinightWails ===
cd /d %~dp0
cd frontend && call npx vite build && cd ..
go build -ldflags "-H windowsgui" -o ArkinightWails.exe .
echo === 编译完成 ===
pause
