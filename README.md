# TrayApp

Windows Tray App for mieru proxy.

# User Guide

1. Copy and paste your mieru client config file to `config.json` in this directory.
2. Double click to run `trayapp.exe`.

You can exit the app by right clicking the icon on windows tray and clicking on exit.

To add the program to windows startup, create a shortcut and copy it to this directory:

```
C:\ProgramData\Microsoft\Windows\Start Menu\Programs\StartUp
```

# Compile

We ship binaries that are ready to use. You are welcomed to compile or bring the binaries by yourself.

### TrayApp

1. Remove old exe files.
2. Install golang.
3. Run `build.bat`.

### mieru

1. Download the latest mieru windows client from https://github.com/enfein/mieru/releases.
2. Rename the executable to `mieru.exe`.

### rsrc

`rsrc.exe` is used to set app icon.

1. Download the latest version from https://github.com/akavel/rsrc/releases.
2. Rename the executable to `rsrc.exe`.

To set your own tray icon, rename your .ico file to `icon.ico` and place it in this directory.

You can set the app icon by running this command:

```
rsrc -manifest main.manifest -ico icon.ico
```
