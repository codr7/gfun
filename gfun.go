package main

import (
	"flag"
	"fmt"
	"github.com/codr7/gfun/lib"
	"log"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()

	var m gfun.M
	m.Init()

	if len(args) == 0 {
		fmt.Printf("Welcome to GFun v%v!\n\n", gfun.Version)

		if err := m.Repl(gfun.DefaultReaders(), os.Stdin, os.Stdout); err != nil {
			log.Fatal(err)
		}
	} else {
		for _, p := range args {
			if err := m.Include(p); err != nil {
				log.Fatal(err)
			}
		}

		m.EmitStop()

		if err := m.Eval(0); err != nil {
			log.Fatal(err)
		}
	}
}

