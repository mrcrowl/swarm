package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
)

// ExitIFError checks if err is non-nil and then shows the message and exits the program
func ExitIfError(err error, message string, args ...interface{}) {
	if err != nil {
		log.Fatalf(message, args...)
		os.Exit(1)
	}
}

// WaitForCtrlC sleeps execution until the program receives a Ctrl+C
func WaitForCtrlC() {
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}

// OpenBrowser opens a URL in the default browser in an OS-specific way
func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Println("Failed to open browser")
	}
}
