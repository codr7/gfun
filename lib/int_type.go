package gfun

import (
	"fmt"
)

type IntType struct {
	BasicType
}

func (self *IntType) Name() *Sym {
	return self.m.Sym("Int")
}

func (self *IntType) AddVal(l *Val, r Val) error {
	if rt := r.Type(); rt != &self.m.IntType {
		return fmt.Errorf("Int add not supported: %v", rt.Name())
	}

	var lv interface{}
	var err error

	if lv, err = l.Data(); err != nil {
		return err
	}

	var rv interface{}
	
	if rv, err = r.Data(); err != nil {
		return err
	}

	l.Init(l.Type(), lv.(int)+rv.(int))
	return nil
}

func (self *IntType) SubVal(l *Val, r Val) error {
	if rt := r.Type(); rt != &self.m.IntType {
		return fmt.Errorf("Int add not supported: %v", rt.Name())
	}

	var lv interface{}
	var err error

	if lv, err = l.Data(); err != nil {
		return err
	}

	var rv interface{}
	
	if rv, err = r.Data(); err != nil {
		return err
	}

	l.Init(l.Type(), lv.(int)-rv.(int))
	return nil
}
