// +build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/getlantern/systray"
)

var (
	cmd       *exec.Cmd
	cmdCancel chan bool
	dir       string
)

const (
	SW_HIDE = 0
)

func main() {
	cmdCancel = make(chan bool)

	// Get the current directory
	var err error
	dir, err = os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// Run the commands when the program starts
	runCommands()

	// Initialize system tray
	systray.Run(onReady, onExit)
}

func onReady() {
	// Set the icon from the ICO file
	iconData := getTrayIcon()
	systray.SetIcon(iconData)
	systray.SetTooltip("Mieru is running")

	// Create an "Exit" menu item
	exitItem := systray.AddMenuItem("Exit", "Exit the program")

	// Goroutine to handle menu item clicks
	go func() {
		for {
			select {
			case <-exitItem.ClickedCh:
				stopCommand()
				exitProgram()
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
	cmd1 := exec.Command("./mieru.exe", "apply", "config", "config.json")
	cmd1.Dir = dir

	// Hide the console window for the Go program
	hideConsoleWindow(cmd1)

	err := cmd1.Start()
	if err != nil {
		fmt.Println("Error running command 1:", err)
	}

	// Build and run the second command: mieru.exe start
	cmd2 := exec.Command("./mieru.exe", "start")
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
	// Read the icon data from an ICO file in the current directory
	iconPath := "icon.ico"
	iconData, err := os.ReadFile(iconPath)
	if err != nil {
		fmt.Println("Error reading icon file:", err)
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
	cmdStop := exec.Command("./mieru.exe", "stop")
	cmdStop.Dir = dir

	// Hide the console window for the Go program
	hideConsoleWindow(cmdStop)

	err := cmdStop.Run()
	if err != nil {
		fmt.Println("Error running stop command:", err)
	}
}
