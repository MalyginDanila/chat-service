package buildinfo

import "runtime/debug"

// BuildInfo заполняется Go автоматически: версия Go, модуль, зависимости, VCS-ревизия.
var BuildInfo *debug.BuildInfo

func init() {
	if bi, ok := debug.ReadBuildInfo(); ok {
		BuildInfo = bi
	}
}
