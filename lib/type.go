package gfun

import (
	"fmt"
	"io"
	"log"
)

type TypeId int

type Type interface {
	Id() TypeId
	Name() *Sym
	Parents() []Type
	Isa(Type) bool
	BoolVal(Val) bool
	EmitVal(Val, Reg, *M) error
	EmitValCall(Val, []Form, Pos, *M) error
	EqVal(Val, Val) bool
	DumpVal(Val, io.Writer)
	String() string
}

type BasicType struct {
	id TypeId
	name *Sym
	parents map[Type]Type
}

func (self *BasicType) Init(m *M, name *Sym, parents...Type) {
	self.id = TypeId(m.nextTypeId)
	m.nextTypeId++
	
	self.name = name
	self.parents = make(map[Type]Type)
	
	for _, p := range parents {
		self.parents[p] = p
		
		for _, pp := range p.Parents() {
			self.parents[pp] = p
		}
	}
}

func (self *BasicType) Id() TypeId {
	return self.id
}

func (self *BasicType) Name() *Sym {
	return self.name
}

func (self *BasicType) Parents() []Type {
	out := make([]Type, len(self.parents))
	i := 0
	
	for _, p := range self.parents {
		out[i] = p
		i++
	}

	return out
}

func (self *BasicType) Isa(parent Type) bool {
	return self.parents[parent] != nil
}

func (self *BasicType) BoolVal(val Val) bool {
	log.Fatalf("Val has no boolean rep: %v", self.name)
	return false
}

func (self *BasicType) EmitVal(val Val, reg Reg, m *M) error {
	return fmt.Errorf("Emit not supported: %v", self)
}

func (self *BasicType) EmitValCall(val Val, args []Form, pos Pos, m *M) error {
	return fmt.Errorf("Call not supported: %v", self)
}

func (self *BasicType) EqVal(l Val, r Val) bool {
	return l.Data() == r.Data()
}

func (self *BasicType) DumpVal(val Val, out io.Writer) {
	fmt.Fprintf(out, "%v", val.Data())
}

func (self *BasicType) String() string {
	return self.name.name
}

func (self *M) BindType(_type Type) {
	n := _type.Name()
	
	if v := self.Env().SetVal(n, false); v == nil {
		log.Fatalf("Dup id: %v", n)
	} else {
		v.Init(&self.MetaType, _type)
	}

	self.types[_type.Id()] = _type
}

func (self *M) GetType(name *Sym) (Type, error) {
	var err error
	var v *Val
	
	if v, err = self.Env().GetVal(name); err != nil {
		return nil, err
	}

	return v.Data().(Type), nil
}


/* Bool */

type BoolType struct {
	BasicType
}

func (self *BoolType) BoolVal(val Val) bool {
	return val.Data().(bool)
}

func (self *BoolType) EmitVal(val Val, reg Reg, m *M) error {
	m.EmitLoadBool(reg, val.Data().(bool))
	return nil
}

func (self *BoolType) DumpVal(val Val, out io.Writer) {
	if (val.Data().(bool)) {
		fmt.Fprintf(out, "T")
	} else {
		fmt.Fprintf(out, "F")
	}
}

/* Fun */

type FunType struct {
	BasicType
}

func (self *FunType) BoolVal(val Val) bool {
	return true
}

func (self *FunType) EmitVal(val Val, reg Reg, m *M) error {
	m.EmitLoadFun(reg, val.Data().(*Fun))
	return nil
}

func (self *FunType) EmitValCall(val Val, args []Form, pos Pos, m *M) error {
	f := val.Data().(*Fun)
	m.EmitEnvPush()

	if len(args) < f.argCount {
		return fmt.Errorf("Missing args for %v: %v %v", f, f.argCount, args)
	}
	
	for i := 0; i < f.argCount; i++ {
		a := args[i]
		
		if err := a.Emit(Reg(i+1), m); err != nil {
			return err
		}
	}

	m.EmitCallI(f)
	return nil
}

func (self *FunType) DumpVal(val Val, out io.Writer) {
	fmt.Fprintf(out, "(Fun %v)", val.Data().(*Fun).name)
}

/* Int */

type IntType struct {
	BasicType
}

func (self *IntType) BoolVal(val Val) bool {
	return val.Data().(int) != 0
}

func (self *IntType) EmitVal(val Val, reg Reg, m *M) error {
	m.EmitLoadInt(reg, val.Data().(int))
	return nil
}

/* Macro */

type MacroType struct {
	BasicType
}

func (self *MacroType) BoolVal(val Val) bool {
	return true
}

func (self *MacroType) EmitVal(val Val, reg Reg, m *M) error {
	m.EmitLoadMacro(reg, val.Data().(*Macro))
	return nil
}

func (self *MacroType) EmitValCall(val Val, args []Form, pos Pos, m *M) error {
	return val.Data().(*Macro).Expand(args, pos, m)
}

func (self *MacroType) DumpVal(val Val, out io.Writer) {
	fmt.Fprintf(out, "(Macro %v)", val.Data().(*Macro).name)
}

/* Meta */

type MetaType struct {
	BasicType
}

func (self *MetaType) BoolVal(val Val) bool {
	return true
}

func (self *MetaType) EmitVal(val Val, reg Reg, m *M) error {
	m.EmitLoadType(reg, val.Data().(Type))
	return nil
}

func (self *MetaType) DumpVal(val Val, out io.Writer) {
	fmt.Fprintf(out, "%v", val.Data().(Type).Name())
}

/* Nil */

type NilType struct {
	BasicType
}

func (self *NilType) BoolVal(val Val) bool {
	return false
}

func (self *NilType) DumpVal(val Val, out io.Writer) {
	fmt.Fprintf(out, "_")
}

/* Var */

type VarType struct {
	BasicType
}

func (self *VarType) EmitVal(val Val, reg Reg, m *M) error {
	m.EmitCopy(reg, val.Data().(Reg))
	return nil
}

func (self *VarType) DumpVal(val Val, out io.Writer) {
	fmt.Fprintf(out, "(Var %v)", val.Data().(Reg))
}
