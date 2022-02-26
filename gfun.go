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

	f, err := m.GetFun(m.Sym("+"))

	if err != nil {
		log.Fatal(err)
	}

	targetReg := m.Env().AllocReg()
	m.RootEnv.Regs[targetReg].Init(&m.FunType, f)
	
	m.EmitCall(targetReg, gfun.CallFlags{})
	m.EmitStop()

	if err := m.Eval(0); err != nil {
		log.Fatal(err)
	}

	v, _ := m.RootEnv.Regs[1].Data()
	fmt.Printf("%v\n", v)

	//fmt.Printf("max: %v\n", 2 & ((1 << gfun.OpCodeBits)-1))
}

