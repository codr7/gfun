package gfun

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unsafe"
)

const (
	OpCount = 1 << OpPcBits
	TypeCount = 1 << OpTypeIdBits
	EnvCount = 1024
	FrameCount = 1024
	SymCount = 1024
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
	
	syms [SymCount]Sym
	symLookup map[string]*Sym
	nextSymId SymId
	
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
	self.symLookup = make(map[string]*Sym)
	self.BeginEnv(nil)

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
			m.debug = !m.debug
			m.Env().Regs[0].Init(&m.BoolType, m.debug)
			return ret, nil
		})

	self.BindNewMacro(self.Sym("dec"), 1,
		func(macro *Macro, args []Form, pos Pos, m *M) error {
			reg, err := m.Env().GetReg(args[0].(*IdForm).id)

			if err != nil {
				return err
			}

			d := 1
			
			if len(args) > 1 {
				d = args[1].(*LitForm).val.Data().(int)
			}	
				
			m.EmitDec(reg, d)
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
	
	self.BindNewFun(self.Sym("dump"), NewArgs().Add(self.Sym("val"), &self.AnyType), nil,
		func(fun *Fun, ret PC, m *M) (PC, error) {
			v := m.Env().Regs[1]
			v.Type().DumpVal(v, os.Stdout)
			fmt.Fprintf(os.Stdout, "\n")
			return ret, nil
		})

	self.BindNewMacro(self.Sym("fun"), 3,
		func(macro *Macro, args []Form, pos Pos, m *M) error {
			var err error
			argForms := args[0].(*SliceForm).items
			var funArgs Args
			
			for i := 0; i < len(argForms); i++ {
				aid := argForms[i].(*IdForm).id
				i++
				atid := argForms[i].(*IdForm).id
				at, err := m.GetType(atid)

				if err != nil {
					return err
				}

				funArgs = funArgs.Add(aid, at)
			}
			
			retId := args[1].(*IdForm).id
			var ret Type

			if ret, err = m.GetType(retId); err != nil {
				return err
			}


			fun := NewFun(nil, funArgs, ret, nil)
			fun.name = m.GenSym(fmt.Sprintf("0x%v", uintptr(unsafe.Pointer(fun))))
			body := NewDoForm(pos, args[2:])			

			if err := fun.Emit(body, m); err != nil {
				return err
			}
			m.Bind(fun.name).Init(&m.FunType, fun)
			m.EmitLoadFun(0, fun)
			return nil
		})

	self.BindNewMacro(self.Sym("fun:"), 4,
		func(macro *Macro, args []Form, pos Pos, m *M) error {
			var err error
			id := args[0].(*IdForm).id
			argForms := args[1].(*SliceForm).items
			var funArgs Args
			
			for i := 0; i < len(argForms); i++ {
				aid := argForms[i].(*IdForm).id
				i++
				atid := argForms[i].(*IdForm).id
				at, err := m.GetType(atid)

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
			body := NewDoForm(pos, args[3:])			

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
			m.EmitEnvBeg()

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
						
					if v.Type() == &m.VarType {
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
				m.Emit(stashOps[i])
			}

			m.EmitEnvEnd()
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
						
					if v.Type() == &m.VarType {
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

	self.BindNewMacro(self.Sym("test"), 2,
		func(macro *Macro, args []Form, pos Pos, m *M) error {
			env := m.Env()
			reg := env.AllocReg()
			m.EmitEnvBeg()
			
			if err := args[0].Emit(reg, m); err != nil {
				return err
			}

			op := m.EmitNop()
			
			for _, f := range args[1:] {
				if err := f.Emit(0, m); err != nil {
					return err
				}
			}

			m.EmitEnvEnd()
			m.EmitStop()
			op.InitTest(reg, m.emitPc)
			return nil
		})

	self.BindNewFun(self.Sym("typeof"),
		NewArgs().Add(self.Sym("val"), &self.AnyType),
		&self.MetaType,
		func(fun *Fun, ret PC, m *M) (PC, error) {
			env := m.Env()
			env.Regs[0].Init(&m.MetaType, env.Regs[1].Type())
			return ret, nil
		})

	self.BindNewFun(self.Sym("+"),
		NewArgs().
			Add(self.Sym("l"), &self.IntType).
			Add(self.Sym("r"), &self.IntType),
		&self.IntType,
		func(fun *Fun, ret PC, m *M) (PC, error) {
			env := m.Env()
			l := env.Regs[1].Data().(int)
			r := env.Regs[2].Data().(int)
			env.Regs[0].Init(&m.IntType, l+r)
			return ret, nil
		})

	self.BindNewFun(self.Sym("-"),
		NewArgs().
			Add(self.Sym("l"), &self.IntType).
			Add(self.Sym("r"), &self.IntType),
		&self.IntType,
		func(fun *Fun, ret PC, m *M) (PC, error) {
			env := m.Env()
			l := env.Regs[1].Data().(int)
			r := env.Regs[2].Data().(int)
			env.Regs[0].Init(&m.IntType, l-r)
			return ret, nil
		})

	self.BindNewFun(self.Sym("<"),
		NewArgs().
			Add(self.Sym("l"), &self.IntType).
			Add(self.Sym("r"), &self.IntType),
		&self.BoolType,
		func(fun *Fun, ret PC, m *M) (PC, error) {
			env := m.Env()
			l := env.Regs[1].Data().(int)
			r := env.Regs[2].Data().(int)
			env.Regs[0].Init(&m.BoolType, l < r)
			return ret, nil
		})
}

func (self *M) Bind(name *Sym) *Val {
	return self.Env().SetVal(name, false)
}

func (self *M) Include(path string) error {
	f, err := os.Open(path)

	if err != nil {
		return fmt.Errorf("Failed opening file: %v %v", path, err)
	}
	
	bin := bufio.NewReader(f)
	pos := NewPos(path, 0, 0)

	for {
		if f, err := ReadForm(defaultReaders, bin, &pos, self); err == io.EOF {
			break
		} else if err != nil {
			return err
		} else if f == nil {
			break
		} else {
			if err := f.Emit(0, self); err != nil {
				return err
			}
		}
	}

	return nil
}
