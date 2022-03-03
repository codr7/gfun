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

func (self *Fun) FuseTailCall(startPc PC, m *M) {
	done := false
	
	for i := m.emitPc-1; !done && i >= startPc; i-- {
		op := &m.ops[i]
		
		switch op.OpCode() {
		case CALLI1:
			if op.CallI1Target() == self {
				op.InitRec()
				log.Printf("Fused tail call at %v", i)
			}

			done = true
		case GOTO, LOAD_NIL, NOP, RET:
			break
		default:
			done = true
		}
	}
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
	
	if err := body.Emit(0, m); err != nil {
		return err
	}

	m.EmitRet()
	skip.InitGoto(m.emitPc)
	startPc = m.Fuse(startPc)
	self.FuseTailCall(startPc, m)
	
	self.body = func(fun *Fun, ret PC, m *M) (PC, error) {
		m.BeginFrame(fun, startPc, ret)
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

	return v.Data().(*Fun), nil
}
