# TrayApp
Windows Tray App for mieru proxy  

# User Guide
Copy and paste your mieru client config file to config.json  
Run trayapp.exe  
You can exit the app by right clicking the icon on windows tray and clicking on exit.  
  
Adding program to windows startup: Create a shortcut and copy it to this directory:  
C:\ProgramData\Microsoft\Windows\Start Menu\Programs\StartUp  

# Compile
Install golang and run build.bat  

Download the latest mieru windows client from here:  
https://github.com/enfein/mieru/releases  
rename the executable to mieru.exe  

rsrc.exe can be used to set app icon. You can download the latest version from here:  
https://github.com/akavel/rsrc/releases  
rename the executable to rsrc.exe  

To set your own tray icon, rename your .ico file to icon.ico and place it in directory.  
You can set the app icon by running this command:  
rsrc -manifest main.manifest -ico icon.ico  
