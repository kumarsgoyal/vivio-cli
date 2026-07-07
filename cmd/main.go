package main

import commands "github.com/viviotv/vivio/cmd/commands"

var (
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

func main() {
	commands.SetVersion(version, commit, buildDate)
	commands.Execute()
}
