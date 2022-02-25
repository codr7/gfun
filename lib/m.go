package gfun

import (
	"sync"
)

type M struct {
	IntType IntType
	nextSym uint64
	syms sync.Map
	rootEnv Env
}

func (self *M) Init() {
	self.IntType.Init(self)
}
