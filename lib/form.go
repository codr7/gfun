package gfun

import (
	"fmt"
)

type Form interface {
	Emit(Reg, *M) error
	EmitCall([]Form, Pos, *M) error
}

type BasicForm struct {
	pos Pos
}

func (self *BasicForm) Init(pos Pos) {
	self.pos = pos
}

func (self *BasicForm) Emit(reg Reg, m *M) error {
	return fmt.Errorf("Emit not supported: %v", self)
}

func (self *BasicForm) EmitCall(args []Form, pos Pos, m *M) error {
	return fmt.Errorf("Call not supported: %v", self)
}

/* Call */

type CallForm struct {
	BasicForm
	target Form
	args []Form
}

func NewCallForm(pos Pos, target Form, args []Form) *CallForm {
	return new(CallForm).Init(pos, target, args)
}

func (self *CallForm) Init(pos Pos, target Form, args []Form) *CallForm {
	self.BasicForm.Init(pos)
	self.target = target
	self.args = args
	return self
}

func (self *CallForm) Emit(reg Reg, m *M) error {
	if  err := self.target.EmitCall(self.args, self.pos, m); err != nil {
		return err
	}

	if reg != 0 {
		m.EmitCopy(reg, 0)
	}
	
	return nil
}

func (self *CallForm) String() string {
	return "()"
}

/* Do */

type DoForm struct {
	BasicForm
	forms []Form
}

func NewDoForm(pos Pos, forms []Form) *DoForm {
	return new(DoForm).Init(pos, forms)
}

func (self *DoForm) Init(pos Pos, forms []Form) *DoForm {
	self.BasicForm.Init(pos)
	self.forms = forms
	return self
}

func (self *DoForm) Emit(reg Reg, m *M) error {
	for _, f := range self.forms {
		if err := f.Emit(reg, m); err != nil {
			return err
		}
	}
	
	return nil
}

func (self *DoForm) String() string {
	return "(do)"
}

/* Id */

type IdForm struct {
	BasicForm
	id *Sym
}

func NewIdForm(pos Pos, id *Sym) *IdForm {
	return new(IdForm).Init(pos, id)
}

func (self *IdForm) Init(pos Pos, id *Sym) *IdForm {
	self.BasicForm.Init(pos)
	self.id = id
	return self
}

func (self *IdForm) Emit(reg Reg, m *M) error {
	v, err := m.Env().GetVal(self.id)

	if err != nil {
		return err
	}

	return v.Type().EmitVal(*v, reg, m)
}

func (self *IdForm) EmitCall(args []Form, pos Pos, m *M) error {
	v, err := m.Env().GetVal(self.id)

	if err != nil {
		return err
	}

	return v.Type().EmitValCall(*v, args, pos, m)	
}

func (self *IdForm) String() string {
	return self.id.name
}

/* Lit */

type LitForm struct {
	BasicForm
	val Val
}

func NewLitForm(pos Pos, val Val) *LitForm {
	return new(LitForm).Init(pos, val)
}

func (self *LitForm) Init(pos Pos, val Val) *LitForm {
	self.BasicForm.Init(pos)
	self.val = val
	return self
}

func (self *LitForm) Emit(reg Reg, m *M) error {
	return self.val.Type().EmitVal(self.val, reg, m)
}

func (self *LitForm) EmitCall(args []Form, pos Pos, m *M) error {
	return self.val.Type().EmitValCall(self.val, args, pos, m)
}

/* Slice */

type SliceForm struct {
	BasicForm
	items []Form
}

func NewSliceForm(pos Pos, items []Form) *SliceForm {
	return new(SliceForm).Init(pos, items)
}

func (self *SliceForm) Init(pos Pos, items []Form) *SliceForm {
	self.BasicForm.Init(pos)
	self.items = items
	return self
}

func (self *SliceForm) String() string {
	return "[]"
}
