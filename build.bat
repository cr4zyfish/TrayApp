@echo off
go install github.com/tc-hib/go-winres@v0.3.3
go-winres simply --icon icon.png
go build -o TrayApp.exe -ldflags -H=windowsgui
echo Finished! Press any key to continue...
pause >nul
