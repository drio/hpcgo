// vim: set ts=2 noet:
package main

type codeGenerator interface {
	genCode(o *Options)
}

var backends = map[string]codeGenerator{}

func main() {
	processArgs()
	backends[opts.BackEnd].genCode(&opts)
}
