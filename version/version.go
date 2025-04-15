package version

import (
	_ "embed"
	"runtime/debug"
	"strings"
)

//go:generate /bin/sh -c "date -u > build_time"
//go:embed build_time
var buildTime string

func BuildTime() string {
	return strings.TrimSpace(buildTime)
}

//go:generate /bin/sh -c "git describe --long --always --dirty --broken --abbrev=8 > build_git"
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
