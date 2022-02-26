package gfun

import (
	"sync"
)

const (
	OpCount = 1 << OpPcBits
)

type PC int

type M struct {
	BoolType BoolType
	FunType FunType
	IntType IntType
	NilType NilType
	
	RootEnv Env
	
	syms sync.Map
	ops [OpCount]Op
	emitPc PC
	env *Env
	frame *Frame
}

func (self *M) Init() {
	self.RootEnv.Init(nil)
	self.env = &self.RootEnv

	self.BoolType.Init(self)
	self.FunType.Init(self)
	self.IntType.Init(self)
	self.NilType.Init(self)
	
	self.BindNewFun(self.Sym("+"), func(self *Fun, callFlags CallFlags, ret PC) (PC, error) {
		var err error
		var l interface{}
		
		if l, err = self.m.env.Regs[1].Data(); err != nil {
			return -1, err
		}
		
		var r interface{}
		
		if r, err = self.m.env.Regs[2].Data(); err != nil {
			return -1, err
		}
		
		self.m.env.Regs[1].Init(&self.m.IntType, l.(int)+r.(int))
		return ret, nil
	})
}

func (self *M) Env() *Env {
	return self.env
}
