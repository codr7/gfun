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

func TestDec(t *testing.T) {
	var m gfun.M
	m.Init()

	m.Env().Regs[1].Init(&m.IntType, 49)
	m.EmitDec(1, 7)
	m.EmitStop()

	if err := m.Eval(0); err != nil {
		t.Fatal(err)
	}

	v := m.Env().Regs[1].Data()

	if v.(int) != 42 {
		t.Fatalf("Wrong result: %v", v)
	}
}

func TestAdd(t *testing.T) {
	var m gfun.M
	m.Init()

	m.Env().Regs[1].Init(&m.IntType, 35)
	m.EmitLoadInt(2, 7)

	f, err := m.GetFun(m.Sym("+"))

	if err != nil {
		t.Fatal(err)
	}

	m.EmitEnvPush()
	m.EmitCallI(f)
	m.EmitStop()

	if err := m.Eval(0); err != nil {
		t.Fatal(err)
	}

	v := m.Env().Regs[0].Data()

	if v.(int) != 42 {
		t.Fatalf("Wrong result: %v", v)
	}
}
