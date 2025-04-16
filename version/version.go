package version

import (
	"runtime/debug"
)

var (
	Commit string
	BuildTime string
)

func CompillerVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.GoVersion
	}
	return ""
}
