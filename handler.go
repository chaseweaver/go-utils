package main

import (
	"flag"
)

type (
	// Command handles individual commands
	Command struct {
		Name        string
		Description string
		Action      func()
		FlagSet     *flag.FlagSet
	}
)

var (
	// Commands handles a map of individual commands
	Commands = make(map[string]Command)
)
