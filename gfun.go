package main

import (
	"fmt"
	"github.com/codr7/gfun/lib"
	"log"
)

func main() {
	var m gfun.M
	m.Init()

	m.RootEnv.Regs[1].Init(&m.IntType, 34)
	m.EmitInc(1)
	m.RootEnv.Regs[2].Init(&m.IntType, 7)
	m.EmitAdd(1, 2)
	m.EmitStop()

	if err := m.Eval(0); err != nil {
		log.Fatal(err)
	}

	v, _ := m.RootEnv.Regs[1].Data()
	fmt.Printf("%v\n", v)
}
