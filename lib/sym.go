package gfun

import (
	"fmt"
)

type SymId int

type Sym struct {
	id SymId
	name string
}

func (self *Sym) Init(id SymId, name string) *Sym {
	self.id = id
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

	if found, _ := self.symLookup[name]; found != nil {
		return found
	}
	
	id := self.nextSymId
	self.nextSymId++
	s := &self.syms[id];
	s.Init(id, name)
	self.symLookup[name] = s
	return s
}

func (self *M) GenSym(name string) *Sym {
	id := self.nextSymId
	self.nextSymId++
	s := &self.syms[id]
	s.Init(id, name)
	return s
}
