package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
	"unicode"

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

type CliTimeout struct {
	Value time.Duration
}

func (c *CliTimeout) UnmarshalText(b []byte) error {
	number := true
	text := string(b)

	for _, char := range text {
		if !unicode.IsDigit(char) {
			number = false

			break
		}
	}

	if number {
		value, err := strconv.ParseUint(text, 0, 64)
		if err != nil {
			return fmt.Errorf("incorrect count of seconds: %w", err)
		}

		c.Value = time.Duration(value) * time.Second

		return nil
	}

	value, err := time.ParseDuration(text)
	if err != nil {
		return fmt.Errorf("incorrect duration: %w", err)
	}

	if value < 0 {
		return fmt.Errorf("duration %s should be >=0", text)
	}

	c.Value = value

	return nil
}

var CLI struct {
	Debug   bool             `short:"d" env:"CHORE_DEBUG" help:"Run in debug mode."`
	Version kong.VersionFlag `short:"V" help:"Show version."`

	List CliCmdList `cmd:"" help:"List namespaces and scripts. Empty namespace lists namespaces."`
	Show CliCmdShow `cmd:"" help:"Show details on a given script"`
	Run  CliCmdRun  `cmd:"" help:"Run chore script."`
}
