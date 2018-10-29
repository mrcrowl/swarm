package ui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"swarm/config"
)

const windowsExecName = "swarm"
const linuxExecName = "./swarm"
const windowsPrompt = "C:\\>"
const linuxPrompt = "$ "

// ChooseBuild chooses the build to use from a list of builds, invoking a user prompt if necessary
func ChooseBuild(builds map[string]*config.RuntimeConfig) *config.RuntimeConfig {
	if len(os.Args) > 1 {
		// build specified as first command line argument
		buildName := os.Args[1]
		if build, found := builds[buildName]; found {
			return build
		}
	}

	var selectedBuild *config.RuntimeConfig
	switch len(builds) {
	case 0: // no builds
		fmt.Printf("No builds found")
		os.Exit(1)

	case 1: // single build
		for k := range builds {
			selectedBuild = builds[k]
			break
		}

	default: // choose from menu
		selectedBuild = chooseBuildFromMenu(builds)
	}

	return selectedBuild
}

// chooseBuildFromMenu presents a menu to select a build
func chooseBuildFromMenu(builds map[string]*config.RuntimeConfig) *config.RuntimeConfig {
	buildNames := enumerateBuildNames(builds)
	fmt.Println("     Hint: use build arg to skip this menu")
	fmt.Printf("              e.g. %s%s %s\n", executablePrompt(), executableName(), buildNames[0])
	fmt.Println("-----------------------------------------------")
	fmt.Println("Choose your build:")

	for {

		appmap := map[string]*config.RuntimeConfig{}
		longestBuildName := 3
		for i, name := range buildNames {
			fmt.Printf("  %d) %s\n", (i + 1), name)
			appmap[strconv.Itoa(i+1)] = builds[name]
			appmap[name] = builds[name]
			if len(name) > longestBuildName {
				longestBuildName = len(name)
			}
			i++
		}
		fmt.Print("  ")
		fmt.Println(strings.Repeat("-", longestBuildName+3))
		fmt.Print("  >")
		reader := bufio.NewReader(os.Stdin)
		lineBytes, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal("Bad input")
		}
		if build, found := appmap[string(lineBytes)]; found {
			return build
		}
	}
}

// enumerateBuildNames selects the name for each build to a list
func enumerateBuildNames(builds map[string]*config.RuntimeConfig) []string {
	buildNames := make([]string, len(builds))
	i := 0
	for name := range builds {
		buildNames[i] = name
		i++
	}
	sort.Strings(buildNames)
	return buildNames
}

// executablePrompt returns the common prompt pattern seen on the platform
func executablePrompt() string {
	if runtime.GOOS == "windows" {
		return windowsPrompt
	}
	return linuxPrompt
}

// executableName returns the appropriate executable name for the platform
func executableName() string {
	if runtime.GOOS == "windows" {
		return windowsExecName
	}
	return linuxExecName
}
