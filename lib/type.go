package gfun

import (
	"fmt"
)

type Type interface {
	fmt.Stringer
	Name() *Sym
	GetValue(interface{}) (interface{}, error)
	Parents() []Type
	Isa(Type) bool
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

func (self *BasicType) GetValue(in interface{}) (interface{}, error) {
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
