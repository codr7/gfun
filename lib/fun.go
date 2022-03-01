package gfun

import (
	"fmt"
	"log"
)

type FunBody = func(*Fun, PC, *M) (PC, error)

type FunArg struct {
	name *Sym
	_type Type
}

type FunArgs []FunArg

func NewFunArgs() FunArgs {
	return nil
}

func (self FunArgs) Add(name *Sym, _type Type) FunArgs {
	 return append(self, FunArg{name: name, _type: _type})
}

type Fun struct {
	name *Sym
	args [FunArgCount]FunArg
	argCount int
	ret Type
	body FunBody
}

func NewFun(name *Sym, args []FunArg, ret Type, body FunBody) *Fun {
	return new(Fun).Init(name, args, ret, body)
}

func (self *Fun) Init(name *Sym, args []FunArg, ret Type, body FunBody) *Fun {
	self.name = name
	
	for i, a := range args {
		self.args[i] = a
		self.argCount++
	}

	self.ret = ret
	self.body = body
	return self
}

func (self *Fun) Call(ret PC, m *M) (PC, error) {
	return self.body(self, ret, m)
}

func (self *Fun) Emit(body Form, m *M) error {
	env := m.Env()
	skip := m.Emit(0)
	startPc := m.emitPc
	
	for i := 0; i < self.argCount; i++ {
		a := self.args[i]
		reg := env.AllocReg()
		env.SetReg(a.name, reg, true)
		env.Regs[reg].Init(&m.VarType, reg)
		m.EmitCopy(reg, Reg(i+1))
	}
	
	if err := body.Emit(m); err != nil {
		return err
	}

	m.EmitRet()
	skip.InitGoto(m.emitPc)
	
	self.body = func(fun *Fun, ret PC, m *M) (PC, error) {
		m.Call(fun, ret)
		return startPc, nil
	}

	return nil
}

func (self *Fun) String() string {
	return fmt.Sprintf("(Fun %v)", self.name)
}

func (self *M) BindNewFun(name *Sym, args []FunArg, ret Type, body FunBody) *Fun {
	f := NewFun(name, args, ret, body)
	
	if v := self.Env().SetVal(name, false); v == nil {
		log.Fatalf("Dup id: %v", name)
	} else {
		v.Init(&self.FunType, f)
	}

	return f
}

func (self *M) GetFun(name *Sym) (*Fun, error) {
	var err error
	var v *Val
	
	if v, err = self.Env().GetVal(name); err != nil {
		return nil, err
	}

	var f interface{}
	
	if f, err = v.Data(); err != nil {
		return nil, err
	}

	return f.(*Fun), nil
}
