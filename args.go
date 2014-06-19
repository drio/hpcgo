// vim: set ts=2 noet:
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Encapsulates details of a Command (hpcgo <command>)
type Command struct {
	Run         func(cmd *Command, args []string)
	UsageLine   string
	Short       string
	Long        string
	Flag        flag.FlagSet
}

// Name returns the command's name: the first word in the usage line.
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
}
