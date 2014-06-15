// vim: set ts=2 noet:
package main

import (
	"fmt"
	"text/template"
	"os"
)

func runCmd() {
	t := template.Must(template.New("cmd").Parse(pbsCmd))
	err := t.Execute(os.Stdout, opts)
	if err != nil {
		fmt.Println(os.Stderr, "Problems parsing pbsCmd string.", err)
	}
}

const pbsCmd = `
mkdir -p {{.Log_dir}}; \
echo '{{.Cmd}}' | \
qsub -N {{.Name}} \
-q {{.Queue}} \
-d $(pwd) \
-o {{.Log_dir}}/{{.Name}}.o \
-e {{.Log_dir}}/{{.Name}}.e \
-l nodes=1:ppn={{.Cores}},mem={{.Memory}} -V
`
