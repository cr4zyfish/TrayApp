@echo off
go install github.com/tc-hib/go-winres@latest
go-winres simply --icon icon.png
go build -o TrayApp.exe -ldflags -H=windowsgui
echo Finished! Press any key to continue...
pause >nul
