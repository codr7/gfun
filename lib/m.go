package gfun

import (
	"log"
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

	self.Bind(self.Sym("T")).Init(&self.BoolType, true)
	self.Bind(self.Sym("F")).Init(&self.BoolType, false)
	self.Bind(self.Sym("_")).Init(&self.NilType, nil)
	
	self.BindNewFun(self.Sym("+"), func(fun *Fun, callFlags CallFlags, ret PC) (PC, error) {
		var err error
		var l interface{}
		
		if l, err = self.env.Regs[1].Data(); err != nil {
			return -1, err
		}
		
		var r interface{}
		
		if r, err = self.env.Regs[2].Data(); err != nil {
			return -1, err
		}
		
		self.env.Regs[1].Init(&self.IntType, l.(int)+r.(int))
		return ret, nil
	})
}

func (self *M) Env() *Env {
	return self.env
}

func (self *M) Bind(name *Sym) *Val {
	v, err := self.Env().SetVal(name)

	if err != nil {
		log.Fatal(err)
	}

	return v
}
