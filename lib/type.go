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
	EmitVal([]Form, Val) ([]Form, error)
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

func (self *BasicType) EmitVal(in []Form, val Val) ([]Form, error) {
	return nil, fmt.Errorf("Emit not supported: %v", self)
}

type BoolType struct {
	BasicType
}

func (self *BoolType) Name() *Sym {
	return self.m.Sym("Bool")
}

func (self *BoolType) BoolVal(val Val) (bool, error) {
	v, err := val.Data()

	if err != nil {
		return false, err
	}
	
	return v.(bool), nil
}

type FunType struct {
	BasicType
}

func (self *FunType) Name() *Sym {
	return self.m.Sym("Fun")
}

func (self *FunType) BoolVal(val Val) (bool, error) {
	return true, nil
}

type IntType struct {
	BasicType
}

func (self *IntType) Name() *Sym {
	return self.m.Sym("Int")
}

func (self *IntType) BoolVal(val Val) (bool, error) {
	v, err := val.Data()

	if err != nil {
		return false, err
	}
	
	return v.(int) != 0, nil
}

func (self *IntType) EmitVal(in []Form, val Val) ([]Form, error) {
	v, err := val.Data()

	if err != nil {
		return nil, err
	}
	
	self.m.EmitLoadInt(0, v.(int))
	return in, nil
}

type NilType struct {
	BasicType
}

func (self *NilType) Name() *Sym {
	return self.m.Sym("Nil")
}

func (self *NilType) BoolVal(val Val) (bool, error) {
	return false, nil
}
