# TrayApp

Windows Tray App for [mieru](https://github.com/enfein/mieru) proxy.

# User Guide

1. Copy and paste your mieru client config file to `config.json` in this directory.
2. Double click to run `TrayApp.exe`.

You can exit the app by right clicking the icon on windows tray and clicking on exit.

# Compile

We ship binaries that are ready to use. You are welcomed to compile or bring the binaries by yourself.

### TrayApp

1. Remove old exe files.
2. Install golang.
3. Run `build.bat`.

### mieru

1. Download the latest mieru windows client from https://github.com/enfein/mieru/releases.
2. Rename the executable to `mieru.exe`.

### App Icon

To set your own App icon, place your icon.png file in this directory. 
To set your own Tray icon, convert you .ico file to base64 and replace it with old icon in the source code.

