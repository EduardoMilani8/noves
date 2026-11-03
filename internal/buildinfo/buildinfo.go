// Package buildinfo guarda a versão do binário, preenchida pelo -ldflags do
// Makefile e do GoReleaser.
package buildinfo

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

// Preenchidos com -ldflags "-X".
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// Resolved devolve a versão, usando os dados do módulo quando o binário veio de
// go install (sem -ldflags).
func Resolved() (version, commit, date string) {
	version, commit, date = Version, Commit, Date
	if version != "dev" {
		return
	}
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	if v := bi.Main.Version; v != "" && v != "(devel)" {
		version = v
	}
	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			if len(s.Value) >= 7 {
				commit = s.Value[:7]
			}
		case "vcs.time":
			date = s.Value
		}
	}
	return
}

// String formata a versão como em "noves version".
func String() string {
	v, c, d := Resolved()
	return fmt.Sprintf("noves %s (commit %s, %s) %s %s/%s", v, c, d, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
