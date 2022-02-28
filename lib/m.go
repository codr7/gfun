package gfun

import (
	"fmt"
	"os"
)

const (
	OpCount = 1 << OpPcBits
	TypeCount = 1 << OpTypeIdBits
)

type PC int

type M struct {
	AnyType BasicType
	BoolType BoolType
	FunType FunType
	IntType IntType
	MacroType MacroType
	MetaType MetaType
	NilType NilType
	VarType VarType
	
	RootEnv Env
	
	syms map[string]*Sym
	types [TypeCount]Type
	nextTypeId TypeId
	ops [OpCount]Op
	emitPc PC
	env *Env
	frame *Frame
	debug bool
}

func (self *M) Init() {
	self.RootEnv.Init(nil)
	self.syms = make(map[string]*Sym)
	self.env = &self.RootEnv

	self.AnyType.Init(self, self.Sym("Any"))
	self.BoolType.Init(self, self.Sym("Bool"), &self.AnyType)
	self.FunType.Init(self, self.Sym("Fun"), &self.AnyType)
	self.IntType.Init(self, self.Sym("Int"), &self.AnyType)
	self.MacroType.Init(self, self.Sym("Macro"), &self.AnyType)
	self.MetaType.Init(self, self.Sym("Meta"), &self.AnyType)
	self.NilType.Init(self, self.Sym("Nil"), &self.AnyType)
	self.VarType.Init(self, self.Sym("Var"), &self.AnyType)

	self.BindType(&self.AnyType)
	self.BindType(&self.BoolType)
	self.BindType(&self.FunType)
	self.BindType(&self.IntType)
	self.BindType(&self.MacroType)
	self.BindType(&self.MetaType)
	self.BindType(&self.NilType)
	self.BindType(&self.VarType)
	
	self.Bind(self.Sym("T")).Init(&self.BoolType, true)
	self.Bind(self.Sym("F")).Init(&self.BoolType, false)
	self.Bind(self.Sym("_")).Init(&self.NilType, nil)

	self.BindNewMacro(self.Sym("begin"), -1,
		func(macro *Macro, args []Form, pos Pos, m *M) error {
			for _, f := range args {
				if err := f.Emit(m); err != nil {
					return err
				}
			}

			return nil
		})
	
	self.BindNewFun(self.Sym("debug"), nil, &self.BoolType,
		func(fun *Fun, ret PC, m *M) (PC, error) {
			self.debug = !self.debug
			self.env.Regs[0].Init(&m.BoolType, self.debug)
			return ret, nil
		})

	self.BindNewFun(self.Sym("dump"), NewFunArgs().Add(self.Sym("val"), &self.AnyType), nil,
		func(fun *Fun, ret PC, m *M) (PC, error) {
			v := self.env.Regs[1]
			v.Type().DumpVal(v, os.Stdout)
			fmt.Fprintf(os.Stdout, "\n")
			return ret, nil
		})

	self.BindNewMacro(self.Sym("fun:"), 4,
		func(macro *Macro, args []Form, pos Pos, m *M) error {
			var err error
			id := args[0].(*IdForm).id
			argForms := args[1].(*SliceForm).items
			var funArgs FunArgs
			
			for i := 0; i < len(argForms); i++ {
				aid := argForms[i].(*IdForm).id
				i++
				atid := argForms[i].(*IdForm).id
				at, err := self.GetType(atid)

				if err != nil {
					return err
				}

				funArgs = funArgs.Add(aid, at)
			}
			
			retId := args[2].(*IdForm).id
			var ret Type

			if ret, err = m.GetType(retId); err != nil {
				return err
			}
			
			fun := m.BindNewFun(id, funArgs, ret, nil)
			body := args[3]			

			if err := fun.Emit(body, m); err != nil {
				return err
			}

			return nil
		})
	
	self.BindNewMacro(self.Sym("if"), 3,
		func(macro *Macro, args []Form, pos Pos, m *M) error {
			if err := args[0].Emit(m); err != nil {
				return err
			}
			
			branch := m.Emit(0)
			truePc := m.emitPc

			if err := args[1].Emit(m); err != nil {
				return err
			}
			
			skip := m.Emit(0)
			falsePc := m.emitPc
			
			if err := args[2].Emit(m); err != nil {
				return err
			}
			
			skip.InitGoto(m.emitPc)
			branch.InitBranch(0, truePc, falsePc)
			return nil
		})

	self.BindNewFun(self.Sym("+"),
		NewFunArgs().
			Add(self.Sym("l"), &self.IntType).
			Add(self.Sym("r"), &self.IntType),
		&self.IntType,
		func(fun *Fun, ret PC, m *M) (PC, error) {
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
		&self.IntType,
		func(fun *Fun, ret PC, m *M) (PC, error) {
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
		&self.BoolType,
		func(fun *Fun, ret PC, m *M) (PC, error) {
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
	return self.Env().SetVal(name, false)
}
