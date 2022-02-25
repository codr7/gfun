package gfun

import (
	"sync"
)

type M struct {
	IntType IntType
	RootEnv Env
	
	syms sync.Map
}

func (self *M) Init() {
	self.RootEnv.Init(nil)
	self.IntType.Init(self)
}
