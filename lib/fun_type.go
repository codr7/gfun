package gfun

import (
)

type FunType struct {
	BasicType
}

func (self *FunType) Name() *Sym {
	return self.m.Sym("Fun")
}

func (self *FunType) BoolVal(val Val) (bool, error) {
	return true, nil
}
