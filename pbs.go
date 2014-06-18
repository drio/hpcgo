// vim: set ts=2 noet:
package main

import (
	"fmt"
	"os"
	"text/template"
)

type pbsBackend struct {
	template string
}

func init() {
	backends["pbs"] = &pbsBackend{pbsCmd}
}

func (p *pbsBackend) genCode(o *Options) {
	t := template.Must(template.New("cmd").Parse(p.template))
	err := t.Execute(os.Stdout, opts)
	if err != nil {
		fmt.Println(os.Stderr, "Problems parsing pbs template.", err)
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
