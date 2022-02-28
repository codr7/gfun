package gfun

import (
	"fmt"
)

type Sym struct {
	name string
}

func (self *Sym) Init(name string) *Sym {
	self.name = name
	return self
}

func (self *Sym) Name() string {
	return self.name
}

func (self *Sym) String() string {
	return self.name
}

func (self *M) Sym(name string, args ...interface{}) *Sym {
	if len(args) > 0 {
		name = fmt.Sprintf(name, args...)
	}

	if found, _ := self.syms[name]; found != nil {
		return found
	}
	
	s := new(Sym).Init(name)
	self.syms[name] = s
	return s
}
