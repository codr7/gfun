package main

import (
	"fmt"
	"github.com/codr7/gfun/lib"
	"log"
	"os"
)

func main() {
	var m gfun.M
	m.Init()

	fmt.Printf("Welcome to GFun v%v!\n\n", gfun.Version)

	//fmt.Printf("max: %v\n", 2 & ((1 << gfun.OpCodeBits)-1))

	if err := m.Repl(gfun.DefaultReaders(), os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

