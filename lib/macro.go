package gfun

import (
	"fmt"
)

type MacroBody = func(*Macro, []Form, Pos, *M) error

type Macro struct {
	name *Sym
	argCount int
	body MacroBody
}

func NewMacro(name *Sym, argCount int, body MacroBody) *Macro {
	return new(Macro).Init(name, argCount, body)
}

func (self *Macro) Init(name *Sym, argCount int, body MacroBody) *Macro {
	self.name = name
	self.argCount = argCount
	self.body = body
	return self
}

func (self *Macro) Expand(args []Form, pos Pos, m *M) error {
	if self.argCount != -1 && len(args) < self.argCount {
		return fmt.Errorf("Missing args for %v: %v", self, args)
	}
	
	return self.body(self, args, pos, m)
}

func (self *Macro) String() string {
	return fmt.Sprintf("(Macro %v)", self.name)
}

func (self *M) BindNewMacro(name *Sym, argCount int, body MacroBody) *Macro {
	f := NewMacro(name, argCount, body)	
	self.Env().SetVal(name, false).Init(&self.MacroType, f)
	return f
}
