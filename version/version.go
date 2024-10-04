package version

import (
	_ "embed"
	"runtime/debug"
	"strings"
)

//go:generate bash -c "date --utc > build_time"
//go:embed build_time
var buildTime string

func BuildTime() string {
	return strings.TrimSpace(buildTime)
}

//go:generate bash -c "git describe --long --always --dirty --broken --abbrev=8 > build_git"
//go:embed build_git
var buildCommit string

func BuildCommit() string {
	return strings.TrimSpace(buildCommit)
}

func CompillerVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.GoVersion
	}
	return ""
}
