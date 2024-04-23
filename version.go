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
		mod := runningModule(info.Deps[i])

		if mod.Path == moduleName {
			return mod.Version
		}
	}
	return "v0.0.0"
}

func runningModule(mod *debug.Module) *debug.Module {
	if mod.Replace != nil {
		return runningModule(mod.Replace)
	} else {
		return mod
	}
}
