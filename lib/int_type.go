package gfun

import (
)

type IntType struct {
	BasicType
}

func (self *IntType) Name() *Sym {
	return self.m.Sym("Int")
}

func (self *IntType) BoolVal(val Val) (bool, error) {
	v, err := val.Data()

	if err != nil {
		return false, err
	}
	
	return v.(int) != 0, nil
}
