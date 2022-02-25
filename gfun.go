package main

import (
	"fmt"
	"github.com/codr7/gfun/lib"
	"log"
)

func main() {
	var m gfun.M
	m.Init()

	m.RootEnv.Regs[1].Init(&m.IntType, 35)
	m.EmitLoadInt(2, 7)
	m.EmitInc(1, 2)
	m.EmitStop()

	if err := m.Eval(0); err != nil {
		log.Fatal(err)
	}

	v, _ := m.RootEnv.Regs[1].Data()
	fmt.Printf("%v\n", v)

	//fmt.Printf("load1 max: %v\n", 1 << (64 - gfun.OpRegBits - gfun.OpCodeBits - 1))
}
