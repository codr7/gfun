package gfun

import (
	"fmt"
)

type CallFlags struct {
	Drop, Memo, Tail bool
}

type FunBody = func(*Fun, CallFlags, PC) (PC, error)

type Fun struct {
	m *M
	name *Sym
	body FunBody
}

func NewFun(m *M, name *Sym, body FunBody) *Fun {
	return new(Fun).Init(m, name, body)
}

func (self *Fun) Init(m *M, name *Sym, body FunBody) *Fun {
	self.m = m
	self.name = name
	self.body = body
	return self
}

func (self *Fun) Call(flags CallFlags, ret PC) (PC, error) {
	return self.body(self, flags, ret)
}

func (self *Fun) Emit(body Form) (PC, error) {
	startPc := self.m.emitPc

	if err := body.Emit(self.m); err != nil {
		return -1, err
	}

	self.m.EmitRet()
	
	self.body = func(fun *Fun, flags CallFlags, ret PC) (PC, error) {
		self.m.Call(fun, flags, ret)
		return startPc, nil
	}

	return startPc, nil
}

func (self *M) BindNewFun(name *Sym, body FunBody) (*Fun, error) {
	f := NewFun(self, name, body)
	
	if err := self.env.SetVal(name, NewVal(&self.FunType, f)); err != nil {
		return nil, err
	}

	return f, nil

}

func (self *M) GetFun(name *Sym) (*Fun, error) {
	var err error
	var v *Val
	
	if v, err = self.env.GetVal(name); err != nil {
		return nil, err
	}

	if v.Type() != &self.FunType {
		return nil, fmt.Errorf("Expected Fun: %v", v)
	}

	var f interface{}
	
	if f, err = v.Data(); err != nil {
		return nil, err
	}

	return f.(*Fun), nil
}
