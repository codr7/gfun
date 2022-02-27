package gfun

import (
	"fmt"
)

const (
	RegCount = 1 << OpRegBits
	ArgCount = 8
)

type Reg int

type Env struct {	
	Regs [RegCount]Val

	outer *Env
	bindings map[*Sym]Reg
	regCount Reg
}

func (self *Env) Init(outer *Env) {
	self.outer = outer

	if outer != nil {
		self.regCount = outer.regCount
		copy(self.Regs[:], outer.Regs[:outer.regCount])
	} else {
		self.regCount = ArgCount+1
	}
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


func (self *Env) GetReg(key *Sym) (Reg, error) {
	tryOuter := func() (Reg, error) {
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

func (self *Env) SetReg(key *Sym, reg Reg) error {
	if self.bindings == nil {
		self.bindings = make(map[*Sym]Reg)
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

func (self *Env) SetVal(key *Sym) (*Val, error) {
	if v := self.FindVal(key); v != nil {
		return nil, fmt.Errorf("Dup id: %v (%v)", key, v)
	}
	
	reg := self.AllocReg()
	self.SetReg(key, reg)
	return &self.Regs[reg], nil
}

func (self *Env) AllocReg() Reg {
	reg := self.regCount
	self.regCount++
	return reg
}
