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
	GetVal(interface{}) (interface{}, error)
	Parents() []Type
	Isa(Type) bool
	BoolVal(Val) bool
	EmitVal(Val, *M) error
	EmitValCall(Val, []Form, Pos, *M) error
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

func (self *BasicType) GetVal(in interface{}) (interface{}, error) {
	return in, nil
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

func (self *BasicType) EmitVal(val Val, m *M) error {
	return fmt.Errorf("Emit not supported: %v", self)
}

func (self *BasicType) EmitValCall(val Val, args []Form, pos Pos, m *M) error {
	return fmt.Errorf("Call not supported: %v", self)
}

func (self *BasicType) DumpVal(val Val, out io.Writer) {
	v, err := val.Data()

	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Fprintf(out, "%v", v)
}

func (self *BasicType) String() string {
	return self.name.name
}

func (self *M) BindType(_type Type) {
	n := _type.Name()
	
	if v := self.env.SetVal(n, false); v == nil {
		log.Fatalf("Dup id: %v", n)
	} else {
		v.Init(&self.MetaType, _type)
	}

	self.types[_type.Id()] = _type
}

func (self *M) GetType(name *Sym) (Type, error) {
	var err error
	var v *Val
	
	if v, err = self.env.GetVal(name); err != nil {
		return nil, err
	}

	var f interface{}
	
	if f, err = v.Data(); err != nil {
		return nil, err
	}

	return f.(Type), nil
}


/* Bool */

type BoolType struct {
	BasicType
}

func (self *BoolType) BoolVal(val Val) bool {
	v, _ := val.Data()
	return v.(bool)
}

func (self *BoolType) EmitVal(val Val, m *M) error {
	v, err := val.Data()

	if err != nil {
		return err
	}
	
	m.EmitLoadBool(0, v.(bool))
	return nil
}

func (self *BoolType) DumpVal(val Val, out io.Writer) {
	v, err := val.Data()

	if err != nil {
		log.Fatal(err)
	}

	if (v.(bool)) {
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

func (self *FunType) EmitVal(val Val, m *M) error {
	v, err := val.Data()

	if err != nil {
		return err
	}
	
	m.EmitLoadFun(0, v.(*Fun))
	return nil
}

func (self *FunType) EmitValCall(val Val, args []Form, pos Pos, m *M) error {
	fd, err := val.Data()

	if err != nil {
		return err
	}

	f := fd.(*Fun)
	m.EmitEnvPush()
	
	for i := 0; i < f.argCount; i++ {
		a := args[i]
		
		if err := a.Emit(m); err != nil {
			return err
		}

		m.EmitCopy(Reg(i+1), 0)
	}


	reg := m.Env().AllocReg()
	m.EmitLoadFun(reg, f)
	m.EmitCall(reg)
	return nil
}

func (self *FunType) DumpVal(val Val, out io.Writer) {
	f, err := val.Data()

	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Fprintf(out, "(Fun %v)", f.(*Fun).name)
}

/* Int */

type IntType struct {
	BasicType
}

func (self *IntType) BoolVal(val Val) bool {
	v, _ := val.Data()
	return v.(int) != 0
}

func (self *IntType) EmitVal(val Val, m *M) error {
	v, err := val.Data()

	if err != nil {
		return err
	}
	
	m.EmitLoadInt(0, v.(int))
	return nil
}

/* Macro */

type MacroType struct {
	BasicType
}

func (self *MacroType) BoolVal(val Val) bool {
	return true
}

func (self *MacroType) EmitVal(val Val, m *M) error {
	v, err := val.Data()

	if err != nil {
		return err
	}
	
	m.EmitLoadMacro(0, v.(*Macro))
	return nil
}

func (self *MacroType) EmitValCall(val Val, args []Form, pos Pos, m *M) error {
	v, err := val.Data()

	if err != nil {
		return err
	}

	return v.(*Macro).Expand(args, pos, m)
}

func (self *MacroType) DumpVal(val Val, out io.Writer) {
	f, err := val.Data()

	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Fprintf(out, "(Macro %v)", f.(*Macro).name)
}

/* Meta */

type MetaType struct {
	BasicType
}

func (self *MetaType) BoolVal(val Val) bool {
	return true
}

func (self *MetaType) EmitVal(val Val, m *M) error {
	v, err := val.Data()

	if err != nil {
		return err
	}
	
	m.EmitLoadType(0, v.(Type))
	return nil
}

func (self *MetaType) DumpVal(val Val, out io.Writer) {
	v, err := val.Data()

	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Fprintf(out, "%v", v.(Type).Name())
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

func (self *VarType) EmitVal(val Val, m *M) error {
	v, err := val.Data()

	if err != nil {
		return err
	}
	
	m.EmitCopy(0, v.(Reg))
	return nil
}

func (self *VarType) DumpVal(val Val, out io.Writer) {
	v, err := val.Data()

	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Fprintf(out, "(Var %v)", v.(Reg))
}
