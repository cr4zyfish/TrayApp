package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/getlantern/systray"
	"golang.org/x/sys/windows/registry"
)

var (
	cmd         *exec.Cmd
	cmdCancel   chan bool
	dir         string
	startOnBoot bool
)

const (
	SW_HIDE = 0
)

const (
	// RegistryKey is the registry key where the startup information is stored
	RegistryKey = `Software\Microsoft\Windows\CurrentVersion\Run`
	// AppName is the name of the application used in the registry
	AppName = "Mieru"
	// ExecutableName is the name of the executable
	ExecutableName = "mieru.exe"
)

func main() {
	cmdCancel = make(chan bool)

	// Get the current directory
	var err error
	dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// Load start-on-boot setting from registry
	loadStartOnBoot()

	// Run the commands when the program starts
	runCommands()

	// Initialize system tray
	systray.Run(onReady, onExit)
}

func onReady() {
	// Set the icon from base64 data
	iconData := getTrayIcon()
	systray.SetIcon(iconData)
	systray.SetTooltip("Mieru is running")

	// Create main menu items
	optionsMenuItem := systray.AddMenuItem("Options", "")
	startOnBootItem := optionsMenuItem.AddSubMenuItemCheckbox("Start on Boot", "Start the program automatically on system boot", startOnBoot)
	openFolderMenuItem := optionsMenuItem.AddSubMenuItem("Open Folder", "Open the folder where the program is run from")
	exitItem := systray.AddMenuItem("Exit", "Exit the program")

	// Goroutine to handle menu item clicks
	go func() {
		for {
			select {
			case <-exitItem.ClickedCh:
				stopCommand()
				exitProgram()
			case <-startOnBootItem.ClickedCh:
				startOnBoot = !startOnBoot
				saveStartOnBoot()
				updateStartOnBoot(startOnBootItem)
			case <-openFolderMenuItem.ClickedCh:
				openFolder()
			}
		}
	}()
}

func onExit() {
	// Clean up resources
	cmdCancel <- true
	if cmd != nil && cmd.Process != nil {
		cmd.Process.Kill()
	}
}

func runCommands() {
	// Build and run the first command: mieru.exe apply config config.json
	cmd1 := exec.Command(filepath.Join(dir, ExecutableName), "apply", "config", "config.json")
	cmd1.Dir = dir

	// Hide the console window for the Go program
	hideConsoleWindow(cmd1)

	err := cmd1.Start()
	if err != nil {
		fmt.Println("Error running command 1:", err)
	}

	// Build and run the second command: mieru.exe start
	cmd2 := exec.Command(filepath.Join(dir, ExecutableName), "start")
	cmd2.Dir = dir

	// Hide the console window for the Go program
	hideConsoleWindow(cmd2)

	// Run the second command asynchronously
	go func() {
		err := cmd2.Run()
		if err != nil {
			fmt.Println("Error running command 2:", err)
		}
	}()

	// Wait for a short duration (adjust as needed)
	time.Sleep(time.Second)
}

func exitProgram() {
	systray.Quit()
}

func getTrayIcon() []byte {
	// Base64-encoded icon data
	iconBase64 := `AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA`

	// Decode the base64 data
	iconData, err := base64.StdEncoding.DecodeString(iconBase64)
	if err != nil {
		fmt.Println("Error decoding icon data:", err)
		return nil
	}

	return iconData
}

// hideConsoleWindow hides the console window for the Go program on Windows
func hideConsoleWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

func stopCommand() {
	// Build and run the stop command: mieru.exe stop
	cmdStop := exec.Command(filepath.Join(dir, ExecutableName), "stop")
	cmdStop.Dir = dir

	// Hide the console window for the Go program
	hideConsoleWindow(cmdStop)

	err := cmdStop.Run()
	if err != nil {
		fmt.Println("Error running stop command:", err)
	}
}

func openFolder() {
	err := exec.Command("explorer", dir).Start()
	if err != nil {
		fmt.Println("Error opening folder:", err)
	}
}

func loadStartOnBoot() {
	key, err := registry.OpenKey(registry.CURRENT_USER, RegistryKey, registry.READ|registry.WRITE)
	if err != nil {
		fmt.Println("Error opening registry key:", err)
		return
	}
	defer key.Close()

	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return
	}

	val, _, err := key.GetStringValue(AppName)
	if err == nil {
		// Mieru value exists in registry
		if val != exePath {
			// Update registry value if the executable path is different
			err = key.SetStringValue(AppName, exePath)
			if err != nil {
				fmt.Println("Error updating registry value:", err)
			}
		}
		startOnBoot = true
	} else if err == registry.ErrNotExist {
		// Mieru value does not exist in registry, set startOnBoot to false
		startOnBoot = false
	} else {
		fmt.Println("Error reading registry value:", err)
	}
}

func saveStartOnBoot() {
	key, err := registry.OpenKey(registry.CURRENT_USER, RegistryKey, registry.WRITE)
	if err != nil {
		fmt.Println("Error opening registry key:", err)
		return
	}
	defer key.Close()

	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return
	}

	val, _, err := key.GetStringValue(AppName)
	if err != nil || val != exePath {
		if startOnBoot {
			err = key.SetStringValue(AppName, exePath)
			if err != nil {
				fmt.Println("Error writing registry value:", err)
			}
		} else {
			err = key.DeleteValue(AppName)
			if err != nil {
				fmt.Println("Error deleting registry value:", err)
			}
		}
	}
}

func updateStartOnBoot(startOnBootItem *systray.MenuItem) {
	if startOnBoot {
		startOnBootItem.Check()
	} else {
		startOnBootItem.Uncheck()
	}
}
