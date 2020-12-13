package main

import (
	"miolang/internal/machine"
	"miolang/internal/parser"
	"os"
)

func main() {
	filename := os.Args[1]
	fp, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	_, cmds := parser.Parse(fp)
	m := machine.NewMachine(cmds)
	m.Run()
}
