# TrayApp

Windows Tray App for [mieru](https://github.com/enfein/mieru) proxy.  

[Download latest release](https://github.com/cr4zyfish/TrayApp/releases/download/v1.0.0/TrayApp-v1.0.0.zip)


# User Guide

1. Copy and paste your mieru client config file to `config.json` in this directory.
2. Double click to run `TrayApp.exe`.

You can exit the app by right clicking the icon on windows tray and clicking on exit.

# Compile

We ship binaries that are ready to use. You are welcomed to compile or bring the binaries by yourself.

### TrayApp
   
1. Remove old exe files.
2. Install [golang](https://dl.google.com/go/go1.21.5.windows-386.msi).
3. Run `build.bat`.

### mieru

1. Download the latest mieru windows client from [here](https://github.com/enfein/mieru/releases).
2. Rename the executable to `mieru.exe`.

### App Icon

- To set your own App icon, place your icon.png file in this directory. Then compile the program.
- To set your own Tray icon, convert your .ico file to base64 and replace it with old iconBase64 data in the source code. Then compile the program.

