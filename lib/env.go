package gfun

import (
	"fmt"
)

const (
	RegCount = 1 << OpRegBits
)

type Reg int

type Env struct {	
	Regs [RegCount]Val

	outer *Env
	bindings map[*Sym]Reg
	regCount Reg
}

func (self *Env) Init(outer *Env) *Env {
	self.outer = outer
	self.bindings = nil
	
	if outer == nil {
		self.regCount = ArgCount+1
	} else {
		self.regCount = outer.regCount

		for i := Reg(0); i < outer.regCount; i++ {
			self.Regs[i] = outer.Regs[i]
		}
	}

	return self
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

func (self *Env) SetReg(key *Sym, reg Reg, force bool) error {
	if self.bindings == nil {
		self.bindings = make(map[*Sym]Reg)
	} else if !force {
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

func (self *Env) SetVal(key *Sym, force bool) *Val {
	if v := self.FindVal(key); v != nil {
		return v
	}
	
	reg := self.AllocReg()
	self.SetReg(key, reg, force)
	return &self.Regs[reg]
}

func (self *Env) AllocReg() Reg {
	reg := self.regCount
	self.regCount++
	return reg
}

func (self *M) Env() *Env {
	if self.envCount == 0 {
		return nil
	}
	
	return &self.envs[self.envCount-1]
}

func (self *M) BeginEnv(outer *Env) *Env {
	e := &self.envs[self.envCount]
	e.Init(outer)
	self.envCount++
	return e
}

func (self *M) EndEnv() *Env {
	self.envCount--
	e := &self.envs[self.envCount]
	self.Env().Regs[0] = e.Regs[0]
	return e
}
