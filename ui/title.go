package ui

import (
	"fmt"
	"log"
	"os"

	flag "github.com/spf13/pflag"
)

const title = "" +
	" _____      ____ _ _ __ _ __ ___  \n" +
	"/ __\\ \\ /\\ / / _` | '__| '_ ` _ \\  v%s\n" +
	"\\__ \\\\ V  V / (_| | |  | | | | | |\n" +
	"|___/ \\_/\\_/ \\__,_|_|  |_| |_| |_| welcomes you\n" +
	"-----------------------------------------------\n"

// PrintTitle outputs the title
func PrintTitle(version string) {
	log.SetOutput(os.Stdout)
	fmt.Printf(title, version)
}

// CheckHelp checks if the --help flag has been set
func CheckHelp(helpFlag *bool) {
	flag.Parse()

	if helpFlag != nil && *helpFlag == true {
		flag.Usage()
		os.Exit(0)
	}
}
