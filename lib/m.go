package gfun

import (
	"fmt"
	"os"
)

const (
	OpCount = 1 << OpPcBits
	TypeCount = 1 << OpTypeIdBits
	EnvCount = 1024
	FrameCount = 1024
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
	
	syms map[string]*Sym
	types [TypeCount]Type
	nextTypeId TypeId
	ops [OpCount]Op
	emitPc PC
	envs [EnvCount]Env
	envCount int
	frames [FrameCount] Frame
	frameCount int
	debug bool
}

func (self *M) Init() {
	self.syms = make(map[string]*Sym)
	self.PushEnv()

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

	self.BindNewMacro(self.Sym("="), 2,
		func(macro *Macro, args []Form, pos Pos, m *M) error {
			if err := args[0].Emit(1, m); err != nil {
				return err
			}

			if err := args[1].Emit(2, m); err != nil {
				return err
			}

			m.EmitEq(1, 2)
			return nil
		})
	
	self.BindNewMacro(self.Sym("bench"), 2,
		func(macro *Macro, args []Form, pos Pos, m *M) error {
			var err error
			
			if err = args[0].Emit(0, m); err != nil {
				return err
			}

			op := m.Emit(0)
			
			if err = args[1].Emit(0, m); err != nil {
				return err
			}

			m.EmitStop()
			op.InitBench(0, m.emitPc)
			return nil
		})

	self.BindNewFun(self.Sym("debug"), nil, &self.BoolType,
		func(fun *Fun, ret PC, m *M) (PC, error) {
			self.debug = !self.debug
			self.Env().Regs[0].Init(&m.BoolType, self.debug)
			return ret, nil
		})

	self.BindNewMacro(self.Sym("dec"), 1,
		func(macro *Macro, args []Form, pos Pos, m *M) error {
			reg, err := self.Env().GetReg(args[0].(*IdForm).id)

			if err != nil {
				return err
			}

			d := 1
			
			if len(args) > 1 {
				d = args[1].(*LitForm).val.Data().(int)
			}	
				
			self.EmitDec(reg, d)
			return nil
		})

	self.BindNewMacro(self.Sym("do"), 1,
		func(macro *Macro, args []Form, pos Pos, m *M) error {
			for _, f := range args {
				if err := f.Emit(0, m); err != nil {
					return err
				}
			}

			return nil
		})
	
	self.BindNewFun(self.Sym("dump"), NewFunArgs().Add(self.Sym("val"), &self.AnyType), nil,
		func(fun *Fun, ret PC, m *M) (PC, error) {
			v := self.Env().Regs[1]
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
			if err := args[0].Emit(0, m); err != nil {
				return err
			}
			
			branch := m.Emit(0)
			truePc := m.emitPc

			if err := args[1].Emit(0, m); err != nil {
				return err
			}
			
			skip := m.Emit(0)
			falsePc := m.emitPc
			
			if err := args[2].Emit(0, m); err != nil {
				return err
			}
			
			skip.InitGoto(m.emitPc)
			branch.InitBranch(0, truePc, falsePc)
			return nil
		})

	self.BindNewMacro(self.Sym("let"), 2,
		func(macro *Macro, args []Form, pos Pos, m *M) error {
			bsf := args[0].(*SliceForm).items
			var stashOps []Op
			m.EmitEnvPush()

			for i := 0; i < len(bsf); i++ {
				k := bsf[i].(*IdForm).id
				i++
				vf := bsf[i]
				env := m.Env()
				v := env.FindVal(k)
				
				if v == nil {
					reg := env.AllocReg()

					if err := env.SetReg(k, reg, false); err != nil {
						return err
					}

					env.Regs[reg].Init(&m.VarType, reg)
					vf.Emit(reg, m)
				} else {
					stashReg := env.AllocReg()
					reg := Reg(-1)
						
					if v.Type() == &self.VarType {
						reg = v.Data().(Reg)

					} else {
						var err error
						
						if reg, err = env.GetReg(k); err != nil {
							return err
						}
					}
					
					m.EmitCopy(stashReg, reg)
					vf.Emit(reg, m)
					var sop Op
					sop.InitCopy(reg, stashReg)
					stashOps = append(stashOps, sop)
				}
			}
						
			for _, f := range args[1:] {
				f.Emit(0, m)
			}

			for i := len(stashOps)-1; i >= 0; i-- {
				self.Emit(stashOps[i])
			}

			m.EmitEnvPop()
			return nil
		})

		self.BindNewMacro(self.Sym("set"), 2,
		func(macro *Macro, args []Form, pos Pos, m *M) error {
			for i := 0; i < len(args); i++ {
				k := args[i].(*IdForm).id
				i++
				vf := args[i]
				env := m.Env()
				v := env.FindVal(k)
				
				if v == nil {
					reg := env.AllocReg()

					if err := env.SetReg(k, reg, false); err != nil {
						return err
					}

					env.Regs[reg].Init(&m.VarType, reg)
					vf.Emit(reg, m)
				} else {
					reg := Reg(-1)
						
					if v.Type() == &self.VarType {
						reg = v.Data().(Reg)

					} else {
						var err error
						
						if reg, err = env.GetReg(k); err != nil {
							return err
						}
					}
					
					vf.Emit(reg, m)
				}
			}

			return nil
		})

	self.BindNewFun(self.Sym("typeof"),
		NewFunArgs().Add(self.Sym("val"), &self.AnyType),
		&self.MetaType,
		func(fun *Fun, ret PC, m *M) (PC, error) {
			env := self.Env()
			env.Regs[0].Init(&self.MetaType, env.Regs[1].Type())
			return ret, nil
		})

	self.BindNewFun(self.Sym("+"),
		NewFunArgs().
			Add(self.Sym("l"), &self.IntType).
			Add(self.Sym("r"), &self.IntType),
		&self.IntType,
		func(fun *Fun, ret PC, m *M) (PC, error) {
			env := self.Env()
			l := env.Regs[1].Data().(int)
			r := env.Regs[2].Data().(int)
			env.Regs[0].Init(&self.IntType, l+r)
			return ret, nil
		})

	self.BindNewFun(self.Sym("-"),
		NewFunArgs().
			Add(self.Sym("l"), &self.IntType).
			Add(self.Sym("r"), &self.IntType),
		&self.IntType,
		func(fun *Fun, ret PC, m *M) (PC, error) {
			env := self.Env()
			l := env.Regs[1].Data().(int)
			r := env.Regs[2].Data().(int)
			env.Regs[0].Init(&self.IntType, l-r)
			return ret, nil
		})

	self.BindNewFun(self.Sym("<"),
		NewFunArgs().
			Add(self.Sym("l"), &self.IntType).
			Add(self.Sym("r"), &self.IntType),
		&self.BoolType,
		func(fun *Fun, ret PC, m *M) (PC, error) {
			env := self.Env()
			l := env.Regs[1].Data().(int)
			r := env.Regs[2].Data().(int)
			env.Regs[0].Init(&self.BoolType, l < r)
			return ret, nil
		})
}

func (self *M) Bind(name *Sym) *Val {
	return self.Env().SetVal(name, false)
}
