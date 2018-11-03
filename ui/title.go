package ui

import (
	"fmt"
	"log"
	"os"
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
