package gfun

import (
	"fmt"
)

type Type interface {
	Name() *Sym
	GetVal(interface{}) (interface{}, error)
	Parents() []Type
	Isa(Type) bool
	BoolVal(Val) (bool, error)
	EmitVal(Val) error
}

type BasicType struct {
	m *M
	parents map[Type]Type
}

func (self *BasicType) Init(m *M, parents...Type) {
	self.m = m

	for _, p := range parents {
		self.parents[p] = p
		
		for _, pp := range p.Parents() {
			self.parents[pp] = p
		}
	}
}

func (self *BasicType) GetVal(in interface{}) (interface{}, error) {
	return in, nil
}

func (self *BasicType) Parents() []Type {
	out := make([]Type, len(self.parents))
	i := 0
	
	for _, p := range self.parents {
		out[i] = p
		i++
	}

	return out
}

func (self *BasicType) Isa(parent Type) bool {
	return self.parents[parent] != nil
}

func (self *BasicType) EmitVal(val Val) error {
	return fmt.Errorf("Emit not supported: %v", self)
}
