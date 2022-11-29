package main

import "github.com/alecthomas/kong"

var version = "dev"

var CLI struct {
	Debug   bool             `short:"d" env:"CHORE_DEBUG" help:"Run in debug mode."`
	Version kong.VersionFlag `short:"V" help:"Show version."`

	List CliCmdList `cmd:"" help:"List namespaces and scripts. Empty namespace lists namespaces."`
	Show CliCmdShow `cmd:"" help:"Show details on a given script"`
	Run  CliCmdRun  `cmd:"" help:"Run chore script."`
}
