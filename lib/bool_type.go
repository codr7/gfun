package gfun

import (
)

type BoolType struct {
	BasicType
}

func (self *BoolType) Name() *Sym {
	return self.m.Sym("Bool")
}

func (self *BoolType) BoolVal(val Val) (bool, error) {
	v, err := val.Data()

	if err != nil {
		return false, err
	}
	
	return v.(bool), nil
}

