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
}
