package gfun

import (
	"fmt"
	"log"
)

type CallFlags struct {
	Drop, Memo, Tail bool
}

type FunBody = func(*Fun, CallFlags, PC) (PC, error)

type Fun struct {
	name *Sym
	body FunBody
}

func NewFun(name *Sym, body FunBody) *Fun {
	return new(Fun).Init(name, body)
}

func (self *Fun) Init(name *Sym, body FunBody) *Fun {
	self.name = name
	self.body = body
	return self
}

func (self *Fun) Call(flags CallFlags, ret PC) (PC, error) {
	return self.body(self, flags, ret)
}

func (self *Fun) Emit(in []Form, body Form, m *M) (PC, []Form, error) {
	startPc := m.emitPc
	var err error
	
	if in, err = body.Emit(in, m); err != nil {
		return -1, nil, err
	}

	m.EmitRet()
	
	self.body = func(fun *Fun, flags CallFlags, ret PC) (PC, error) {
		m.Call(fun, flags, ret)
		return startPc, nil
	}

	return startPc, in, nil
}

func (self *M) BindNewFun(name *Sym, body FunBody) *Fun {
	f := NewFun(name, body)
	
	if v, err := self.env.SetVal(name); err != nil {
		log.Fatal(err)
	} else {
		v.Init(&self.FunType, f)
	}

	return f

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
