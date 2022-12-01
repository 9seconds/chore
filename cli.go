package main

import "github.com/alecthomas/kong"

var version = "dev"

var CLI struct {
	Debug   bool             `short:"d" env:"CHORE_DEBUG" help:"Run in debug mode."`
	Version kong.VersionFlag `short:"V" help:"Show version."`

	List CliCmdList `cmd:"" aliases:"l" help:"List namespaces and scripts. Empty namespace lists namespaces."`
	Show CliCmdShow `cmd:"" aliases:"s" help:"Show details on a given script."`
	Edit CliCmdEdit `cmd:"" aliases:"e" help:"Edit chore script."`
	Run  CliCmdRun  `cmd:"" aliases:"r" help:"Run chore script."`
}
