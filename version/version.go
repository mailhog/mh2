package version

import "fmt"

const (
	// Version is the version
	Version = "2.0.0-alpha"
	// BuildDate is the build date
	BuildDate = "2017-01-10 22:57"
)

// String returns a formatted version and build date string
func String() string {
	return fmt.Sprintf("%s (%s)", Version, BuildDate)
}
