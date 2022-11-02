package main

import (
	"context"
	"fmt"
	"os"

	"github.com/9seconds/chore/internal/env"
	"github.com/alecthomas/kong"
)

var version = "dev"

type Context struct {
	context.Context
}

type CliNamespace struct {
	Value string
}

func (c *CliNamespace) UnmarshalText(b []byte) error {
	text := string(b)

	if text != "." {
		c.Value = text

		return nil
	}

	text, ok := os.LookupEnv(env.EnvNamespace)
	if !ok {
		return fmt.Errorf("Namespace is dotted but no value for %s is provided", env.EnvNamespace)
	}

	c.Value = text

	return nil
}

func (c CliNamespace) String() string {
	return c.Value
}

func (c CliNamespace) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

var CLI struct {
	Debug   bool             `short:"d" env:"CHORE_DEBUG" help:"Run in debug mode."`
	Version kong.VersionFlag `short:"V" help:"Show version."`

	List CliCmdList `cmd:"" help:"List namespaces and scripts. Empty namespace lists namespaces."`
	Show CliCmdShow `cmd:"" help:"Show details on a given script"`
	Run  CliCmdRun  `cmd:"" help:"Run chore script."`
}
