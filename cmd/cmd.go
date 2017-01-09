package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mailhog/mh2/version"
)

// Main handles default commands which exist across applications
func Main(f func()) {
	if len(os.Args) > 1 {
		switch strings.ToLower(os.Args[1]) {
		case "version":
			fmt.Fprintf(os.Stderr, "%s (%s)\n", version.Version, version.BuildDate)
			os.Exit(0)
			return
		}
	}
	f()
}
