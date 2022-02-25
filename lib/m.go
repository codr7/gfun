package gfun

import (
	"sync"
)

const (
	OpCount = 1024
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
}

func (self *M) Init() {
	self.RootEnv.Init(nil)
	self.BoolType.Init(self)
	self.FunType.Init(self)
	self.IntType.Init(self)
}
