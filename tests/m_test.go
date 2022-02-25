package tests

import (
	"github.com/codr7/gfun/lib"
	"testing"
)

func TestM(t *testing.T) {
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
