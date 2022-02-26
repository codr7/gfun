package tests

import (
	"github.com/codr7/gfun/lib"
	"testing"
)

func TestSym(t *testing.T) {
	var m gfun.M
	m.Init()

	if s := m.Sym("foo"); s.Name() != "foo" {
		t.Fatalf("Wrong name")
	}
	
	if m.Sym("foo") != m.Sym("foo") {
		t.Fatalf("Invalid sym")
	}


	if m.Sym("foo") == m.Sym("bar") {
		t.Fatalf("Dup sym")
	}
}

func TestInc(t *testing.T) {
	var m gfun.M
	m.Init()

	m.RootEnv.Regs[1].Init(&m.IntType, 35)
	m.EmitLoadInt(2, 7)
	m.EmitInc(1, 2)
	m.EmitStop()

	if err := m.Eval(0); err != nil {
		t.Fatal(err)
	}

	if v, err := m.RootEnv.Regs[1].Data(); err != nil {
		t.Fatal(err)
	} else if v.(int) != 42 {
		t.Fatalf("Wrong result: %v", v)
	}
}

func TestAdd(t *testing.T) {
	var m gfun.M
	m.Init()

	m.RootEnv.Regs[1].Init(&m.IntType, 35)
	m.EmitLoadInt(2, 7)

	f, err := m.GetFun(m.Sym("+"))

	if err != nil {
		t.Fatal(err)
	}

	targetReg := m.Env().AllocReg()
	m.RootEnv.Regs[targetReg].Init(&m.FunType, f)
	
	m.EmitCall(targetReg, gfun.CallFlags{})
	m.EmitStop()

	if err := m.Eval(0); err != nil {
		t.Fatal(err)
	}

	if v, err := m.RootEnv.Regs[1].Data(); err != nil {
		t.Fatal(err)
	} else if v.(int) != 42 {
		t.Fatalf("Wrong result: %v", v)
	}
}
