package moovgo

import (
	"runtime/debug"
)

const (
	moduleName = "github.com/moovfinancial/moov-go"
)

func Version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "v0.0.0"
	}
	if info.Main.Path == moduleName {
		return info.Main.Version
	}
	for i := range info.Deps {
		if info.Deps[i].Path == moduleName {
			return info.Deps[i].Version
		}
	}
	return "v0.0.0"
}
