package gfun

import (
	"fmt"
)

const (
	RegCount = 256
	ArgCount = 8
	RetCount = 4
)

type Env struct {	
	Regs [RegCount]Val
	Args [ArgCount]Val
	Rets [RetCount]Val

	outer *Env
	bindings map[*Sym]int
	regCount int
}

func (self *Env) Init(outer *Env) {
	self.outer = outer
}

func (self *Env) FindVal(key *Sym) *Val {
	if self.bindings == nil {
		return nil
	}
	
	reg, ok := self.bindings[key]
	
	if !ok {
		return nil
	}

	return &self.Regs[reg]
}


func (self *Env) GetReg(key *Sym) (int, error) {
	tryOuter := func() (int, error) {
		if self.outer == nil {
			return -1, fmt.Errorf("Unknown id: %v", key)
		}
		
		return self.outer.GetReg(key)
	}
	
	if self.bindings == nil {
		return tryOuter()
	}
	
	reg, ok := self.bindings[key]
	
	if !ok {
		return tryOuter()
	}

	return reg, nil
}

func (self *Env) SetReg(key *Sym, reg int) error {
	if self.bindings == nil {
		self.bindings = make(map[*Sym]int)
	} else {
		if v, dup := self.bindings[key]; dup {
			return fmt.Errorf("Dup id: %v (%v)", key, v)
		}
	}

	self.bindings[key] = reg
	return nil
}

func (self *Env) GetVal(key *Sym) (*Val, error) {
	reg, err := self.GetReg(key)

	if err != nil {
		return nil, err
	}

	return &self.Regs[reg], nil
}

func (self *Env) AllocReg(key *Sym) (int, error) {
	reg := self.regCount
	self.regCount++
	return reg, self.SetReg(key, reg)
}
