@echo off
go build -o trayapp.exe -ldflags -H=windowsgui
rsrc -manifest main.manifest -ico icon.ico
echo Finished! Press any key to continue...
pause >nul
