package ui

import "fmt"

const title = "" +
	" _____      ____ _ _ __ _ __ ___  \n" +
	"/ __\\ \\ /\\ / / _` | '__| '_ ` _ \\\n" +
	"\\__ \\\\ V  V / (_| | |  | | | | | |\n" +
	"|___/ \\_/\\_/ \\__,_|_|  |_| |_| |_| welcomes you\n\n"

func PrintTitle() {
	fmt.Print(title)
}
