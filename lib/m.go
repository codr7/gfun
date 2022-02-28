package gfun

import (
	"log"
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
	
	syms map[string]*Sym
	nextTypeId TypeId
	ops [OpCount]Op
	emitPc PC
	env *Env
	frame *Frame
}

func (self *M) Init() {
	self.RootEnv.Init(nil)
	self.syms = make(map[string]*Sym)
	self.env = &self.RootEnv

	self.BoolType.Init(self, self.Sym("Bool"))
	self.FunType.Init(self, self.Sym("Fun"))
	self.IntType.Init(self, self.Sym("Int"))
	self.NilType.Init(self, self.Sym("Nil"))
	
	self.Bind(self.Sym("T")).Init(&self.BoolType, true)
	self.Bind(self.Sym("F")).Init(&self.BoolType, false)
	self.Bind(self.Sym("_")).Init(&self.NilType, nil)
	
	self.BindNewFun(self.Sym("+"),
		NewFunArgs().
			Add(self.Sym("l"), &self.IntType).
			Add(self.Sym("r"), &self.IntType),
		func(fun *Fun, callFlags CallFlags, ret PC) (PC, error) {
		var err error
		var l interface{}
		
		if l, err = self.env.Regs[1].Data(); err != nil {
			return -1, err
		}
		
		var r interface{}
		
		if r, err = self.env.Regs[2].Data(); err != nil {
			return -1, err
		}
		
		self.env.Regs[0].Init(&self.IntType, l.(int)+r.(int))
		return ret, nil
	})

	self.BindNewFun(self.Sym("-"),
		NewFunArgs().
			Add(self.Sym("l"), &self.IntType).
			Add(self.Sym("r"), &self.IntType),
		func(fun *Fun, callFlags CallFlags, ret PC) (PC, error) {
		var err error
		var l interface{}
		
		if l, err = self.env.Regs[1].Data(); err != nil {
			return -1, err
		}
		
		var r interface{}
		
		if r, err = self.env.Regs[2].Data(); err != nil {
			return -1, err
		}
		
		self.env.Regs[0].Init(&self.IntType, l.(int)-r.(int))
		return ret, nil
	})

	self.BindNewFun(self.Sym("<"),
		NewFunArgs().
			Add(self.Sym("l"), &self.IntType).
			Add(self.Sym("r"), &self.IntType),
		func(fun *Fun, callFlags CallFlags, ret PC) (PC, error) {
		var err error
		var l interface{}
		
		if l, err = self.env.Regs[1].Data(); err != nil {
			return -1, err
		}
		
		var r interface{}
		
		if r, err = self.env.Regs[2].Data(); err != nil {
			return -1, err
		}
		
		self.env.Regs[0].Init(&self.BoolType, l.(int) < r.(int))
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
