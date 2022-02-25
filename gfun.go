package main

import (
	"fmt"
	"github.com/codr7/gfun/lib"
	"log"
)

func main() {
	var m gfun.M
	m.Init()

	m.RootEnv.Regs[42].Init(&m.IntType, 0)
	m.EmitInc(42)
	m.EmitStop()

	if err := m.Eval(0); err != nil {
		log.Fatal(err)
	}

	v, _ := m.RootEnv.Regs[42].Data()
	fmt.Printf("%v\n", v)
}
