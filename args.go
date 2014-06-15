// vim: set ts=2 noet:
package main

import (
	"fmt"
	"flag"
	"os"
	"text/template"
)

type Options struct {
  Name, Memory, Cores, Queue, Log_dir, Deps string
}

var defaults = Options { "", "4Gb", "1", "analysis", "submit_logs", "", }
var opts = Options{}

func processArgs() {
	flag.Usage = usage
	flag.Parse()
	if opts.Name == "" {
		error("Job name not provided.")
	}
}

func init() {
	flag.StringVar(&opts.Name, "s", defaults.Name, "Job name.",)
	flag.StringVar(&opts.Memory, "m", defaults.Memory, "Amount of memory to request.",)
	flag.StringVar(&opts.Cores, "c", defaults.Cores, "Number of cores to request.",)
	flag.StringVar(&opts.Queue, "q", defaults.Queue, "Queue to use.",)
	flag.StringVar(&opts.Log_dir, "l", defaults.Log_dir, "Directory to use for logging.",)
	flag.StringVar(&opts.Deps, "d", defaults.Deps, "List of dependencies.",)
}

func error(msg string) {
	fmt.Fprintln(os.Stderr, "Ups!:", msg)
	usage()
	os.Exit(1)
}

func usage() {
	t := template.Must(template.New("usage").Parse(usageString))
	err := t.Execute(os.Stderr, opts)
	if err != nil {
		fmt.Println(os.Stderr, "Problems parsing usage string.", err)
	}
}

const usageString = `
Usage : hpggo -s job_name [<options>] <cmd>

  -m memory ({{.Memory}})
  -c number of cores ({{.Cores}})
  -q queue to use ({{.Queue}})
  -l directory to use to dump logs ({{.Log_dir}})
  -d list of dependencies. dep1:dep2:.... :depn ("{{.Deps}}")
  -f file to locate list of dependencies. One line, one depency.

  <cmd> is the set of shell commands to execute
  If no <cmd> provided, we will read from stdin

Examples:
  # Iterate over a bunch of files and hgpgo a jobs using each file
  $ ls *.fastq | xargs -i hgpgo -s fmi.{} -m 8G -c 8 "sga index --no-reverse -d 5000000 -t 8 {}"

  # Similar to before but using shell's for
  $ for i in $(ls ../input/reads.*.fastq); do F=$(basename $i .fastq); hgpgo -s pp.$F "sga preprocess -o $F.pp.fastq --pe-mode 2 $i"; done

  # hgpgo a two jobs, the second one has to run after the first one completes
  $ hgpgo -s one "touch ./one.txt" | bash > /tmp/deps.txt ; hgpgo -s two -f /tmp/deps.txt  "sleep 2;touch ./two.txt" | bash

  # Same as before but now we specify the jobid in the command line instead in a file
  $ hgpgo -s filter -d 3678650.sug-moab -m 20G -c 6 "sga fm-merge -m 65 -t 6 final.filter.pass.fa"

  # And my favourite one.
  # The second command reads the jobid from the standard input and uses it as dep
  $ (hgpgo -s one "sleep 15; touch ./one.txt" | bash) | hgpgo -s two -f -  "sleep 2;touch ./two.txt" | bash

`
