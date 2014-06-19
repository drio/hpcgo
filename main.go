// vim: set ts=2 noet:
package main

import (
	"flag"
	"fmt"
	"os"
	"text/template"
	"io"
)

var commands = []*Command{
	cmdSubmit,
}

func main() {
	flag.Usage = mainUsage
	flag.Parse()

	// No arguments
	args := flag.Args()
	if len(args) < 1 {
		mainUsage()
		return
	}

	// User requesting help
	if args[0] == "help" {
		help(args[1:])
		return
	}

	// Is the second argument a valid command?
  // If so, process the rest of the command line and run the command
	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() { cmd.Usage() }
			cmd.Flag.Parse(args[1:])
			args = cmd.Flag.Args()
			cmd.Run(cmd, args)
			os.Exit(2)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "hpcgo: unknown subcommand %q\nRun 'go help' for usage.\n", args[0])
	os.Exit(2)
}

func error(msg string) {
	fmt.Fprintln(os.Stderr, "Ups!:", msg)
	mainUsage()
	os.Exit(1)
}

func help(args []string) {
	if len(args) == 0 {
		mainUsage()
		return
	}

	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: hpcgo help command\n\nToo many arguments given.\n")
		os.Exit(2)
	}

	// User wants help for a particular command
	arg := args[0]
	for _, cmd := range commands {
		if cmd.Name() == arg {
			tmpl(os.Stdout, helpTemplate, cmd)
			return
		}
	}
}

func mainUsage() {
	tmpl(os.Stdout, mainUsageTemplate, commands)
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

var helpTemplate = `Usage: {{.UsageLine}}
{{.Long}}
`

var mainUsageTemplate = `hpcgo is to help you send jobs to an HPC cluster.

Use "go help [command]" for more info about a command.

The commands are:
{{range .}}
    {{.Name | printf "%-11s"}} {{.Short}}
{{end}}
`
