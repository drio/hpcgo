// vim: set ts=2 noet:
package main

import (
	"fmt"
	"strings"
	"bytes"
)

type codeGenerator interface {
	genCode(o *Options)
}

var backends = map[string]codeGenerator{}

var cmdSubmit = &Command{
UsageLine: "submit -s job_name [<options>] <cmd>",
	Short:     "Generate necessary code to submit cluster jobs",
	Long: `
  -b backend cluster to use when generating target code ({{.BackEnd}})
  -m memory ({{.Memory}})
  -c number of cores ({{.Cores}})
  -q queue to use ({{.Queue}})
  -l directory to use to dump logs ({{.Log_dir}})
  -d list of dependencies. dep1:dep2:.... :depn ("{{.Deps}}")
  -f file to locate list of dependencies. One line, one depency.

  <cmd> is the set of shell commands to execute
  If no <cmd> provided, we will read from stdin

Examples:
  # Iterate over a bunch of files and hpcgo a jobs using each file
  $ ls *.fastq | xargs -i hpcgo submit -s fmi.{} -m 8G -c 8 "sga index --no-reverse -d 5000000 -t 8 {}"

  # Similar to before but using shell's for
  $ for i in $(ls ../input/reads.*.fastq); do F=$(basename $i .fastq); hpcgo submit -s pp.$F "sga preprocess -o $F.pp.fastq --pe-mode 2 $i"; done

  # hpcgo a two jobs, the second one has to run after the first one completes
  $ hpcgo submit -s one "touch ./one.txt" | bash > /tmp/deps.txt ; hpcgo submit -s two -f /tmp/deps.txt  "sleep 2;touch ./two.txt" | bash

  # Same as before but now we specify the jobid in the command line instead in a file
  $ hpcgo submit -s filter -d 3678650.sug-moab -m 20G -c 6 "sga fm-merge -m 65 -t 6 final.filter.pass.fa"

  # And my favourite one.
  # The second command reads the jobid from the standard input and uses it as dep
  $ (hpcgo submit -s one "sleep 15; touch ./one.txt" | bash) | hpcgo submit -s two -f -  "sleep 2;touch ./two.txt" | bash
`,
}

type Options struct {
	Name, Memory, Cores, Queue, Log_dir, Deps, Cmd, BackEnd string
}

var opts = Options{}

var defaults = Options{
	"", "4Gb", "1", "analysis", "submit_logs", "", "", "pbs",
}

func setDefaultsInHelp() {
	new_long := bytes.NewBuffer([]byte(""))
	tmpl(new_long, cmdSubmit.Long, defaults)
	cmdSubmit.Long = new_long.String()
}

func init() {
	setDefaultsInHelp()
	cmdSubmit.Run = runSubmit
	addSubmitFlags(cmdSubmit)
}


func addSubmitFlags(cmd *Command) {
	cmd.Flag.StringVar(&opts.Name, "s", defaults.Name, "Job name.")
	cmd.Flag.StringVar(&opts.Memory, "m", defaults.Memory, "Amount of memory to request.")
	cmd.Flag.StringVar(&opts.Cores, "c", defaults.Cores, "Number of cores to request.")
	cmd.Flag.StringVar(&opts.Queue, "q", defaults.Queue, "Queue to use.")
	cmd.Flag.StringVar(&opts.Log_dir, "l", defaults.Log_dir, "Directory to use for logging.")
	cmd.Flag.StringVar(&opts.Deps, "d", defaults.Deps, "List of dependencies.")
	cmd.Flag.StringVar(&opts.BackEnd, "b", defaults.BackEnd, "Backend cluster target to use.")
}

func runSubmit(cmd *Command, args []string) {
	if opts.Name == "" {
		error("Job name not provided.")
	}

	if _, found := backends[opts.BackEnd]; !found {
		error(fmt.Sprintf("<%s> is not a valid backend.", opts.BackEnd))
	}

	opts.Cmd = strings.Join(args, " ")

	// Generate target command for specific backend
	backends[opts.BackEnd].genCode(&opts)
}
