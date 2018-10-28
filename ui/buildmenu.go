package ui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"strconv"
	"swarm/config"
)

const windowsExecName = "swarm.exe"
const linuxExecName = "./swarm"

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
	fmt.Printf("Choose your build: (hint: default argument is the build, e.g. %s %s)\n", executableName(), buildNames[0])

	for {

		appmap := map[string]*config.RuntimeConfig{}
		for i, name := range buildNames {
			fmt.Printf("  %d) %s\n", (i + 1), name)
			appmap[strconv.Itoa(i+1)] = builds[name]
			appmap[name] = builds[name]
			i++
		}
		fmt.Println("  -----------------------")
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

// executableName returns the appropriate executable name for the platform
func executableName() string {
	if runtime.GOOS == "windows" {
		return windowsExecName
	}
	return linuxExecName
}
