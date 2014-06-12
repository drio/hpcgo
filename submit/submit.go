// vim: set ts=4 et:
package main

import (
    "flag"
    "fmt"
    "os"
    //"github.com/drio/submit.go"
)

const HELP string = `
HELP
    here1
    here
`

var validActions = map[string]func(*Options) {
    "submit" : submitAction,
}

type Options struct {
    name, memory, cores, queue, log_dir, deps string
}

var defaults = Options {
    "", "4Gb", "1", "analysis", "submit_logs", "",
}

func submitAction(opts *Options) {
    fmt.Println("In submit Action", opts.memory)
}

func usage(args ...string) {
    exit_code := 0
    if len(args) == 1 {
        fmt.Println("ERROR: ", args[0])
        exit_code = 1
    }
    println(HELP)
    os.Exit(exit_code)
}

func processArgs() *Options {
    pa := Options{}
    flag.StringVar(&pa.memory, "name", defaults.memory, "Job name")
    flag.StringVar(&pa.memory, "memory", defaults.memory, "Memory to request.")
    flag.StringVar(&pa.cores, "cores", defaults.cores, "Num of Cores to request.")
    flag.StringVar(&pa.queue, "queue", defaults.queue, "Queue to use.")
    flag.StringVar(&pa.log_dir, "log_dir", defaults.log_dir, "Log dir output.")
    flag.StringVar(&pa.deps, "deps", defaults.deps, "List of dependencies.")
    flag.Parse()
    return &pa
}

func validateCli(options *Options) string {
    rest_args := flag.Args()

    if len(rest_args) != 1 {
        usage("Wrong number of arguments")
    }

    action := rest_args[0]
    if _, ok := validActions[action]; !ok {
        usage("Invalid action")
    }

    return action
}

func main() {
    args   := processArgs()
    action := validateCli(args)
    validActions[action](args)
}
