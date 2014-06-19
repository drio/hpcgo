// vim: set ts=2 noet:
package main

import (
	"fmt"
)

const VERSION = "0.1.0"

var cmdVersion = &Command{
	Run:       runVersion,
	UsageLine: "version",
	Short:     "print hpcgo version",
	Long:      `Version prints the hpcgo version.`,
}

func runVersion(cmd *Command, args []string) {
	if len(args) != 0 {
		cmd.Usage()
	}

	fmt.Printf("hpcgo version %s\n", VERSION)
}
