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
	BoolVal(Val) (bool, error)
	EmitVal(Val, *M) error
	EmitValCall(Val, []Form, *M) error
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

func (self *BasicType) EmitVal(val Val, m *M) error {
	return fmt.Errorf("Emit not supported: %v", self)
}

func (self *BasicType) EmitValCall(val Val, args []Form, m *M) error {
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

/* Bool */

type BoolType struct {
	BasicType
}

func (self *BoolType) BoolVal(val Val) (bool, error) {
	v, err := val.Data()

	if err != nil {
		return false, err
	}
	
	return v.(bool), nil
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

func (self *FunType) BoolVal(val Val) (bool, error) {
	return true, nil
}

func (self *FunType) EmitVal(val Val, m *M) error {
	v, err := val.Data()

	if err != nil {
		return err
	}
	
	m.EmitLoadFun(0, v.(*Fun))
	return nil
}

func (self *FunType) EmitValCall(val Val, args []Form, m *M) error {
	for i, a := range args {
		if err := a.Emit(m); err != nil {
			return err
		}

		m.EmitMove(Reg(i+1), 0)
	}

	f, err := val.Data()

	if err != nil {
		return err
	}

	reg := m.Env().AllocReg()
	m.EmitLoadFun(reg, f.(*Fun))
	m.EmitCall(reg, CallFlags{})
	return nil
}

func (self *FunType) DumpVal(val Val, out io.Writer) {
	f, err := val.Data()

	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Fprintf(out, "%v()", f.(*Fun).name)
}

/* Int */

type IntType struct {
	BasicType
}

func (self *IntType) BoolVal(val Val) (bool, error) {
	v, err := val.Data()

	if err != nil {
		return false, err
	}
	
	return v.(int) != 0, nil
}

func (self *IntType) EmitVal(val Val, m *M) error {
	v, err := val.Data()

	if err != nil {
		return err
	}
	
	m.EmitLoadInt(0, v.(int))
	return nil
}

/* Nil */

type NilType struct {
	BasicType
}

func (self *NilType) BoolVal(val Val) (bool, error) {
	return false, nil
}
