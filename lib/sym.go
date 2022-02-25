package gfun

import (
	"fmt"
	"sync/atomic"
)

type Sym uint64

func (self *M) Sym(name string, args ...interface{}) Sym {
	var s Sym

	if len(args) > 0 {
		name = fmt.Sprintf(name, args...)
	}

	if found, _ := self.syms.Load(name); found != nil {
		return found.(Sym)
	}
	
	s = self.NextSym()
	
	if ls, found := self.syms.LoadOrStore(name, s); found {
		return ls.(Sym)
	}

	return s 
}

func (self *M) NextSym() Sym {
	return Sym(atomic.AddUint64(&self.nextSym, 1))
}
