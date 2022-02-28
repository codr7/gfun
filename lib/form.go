package gfun

import (
	"fmt"
)

type Form interface {
	Emit(*M) error
	EmitCall([]Form, Pos, *M) error
}

type BasicForm struct {
	pos Pos
}

func (self *BasicForm) Init(pos Pos) {
	self.pos = pos
}

func (self *BasicForm) Emit(m *M) error {
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

func (self *CallForm) Emit(m *M) error {
	if  err := self.target.EmitCall(self.args, self.pos, m); err != nil {
		return err
	}

	return nil
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

func (self *IdForm) Emit(m *M) error {
	v, err := m.env.GetVal(self.id)

	if err != nil {
		return err
	}

	return v.Type().EmitVal(*v, m)
}

func (self *IdForm) EmitCall(args []Form, pos Pos, m *M) error {
	v, err := m.env.GetVal(self.id)

	if err != nil {
		return err
	}

	return v.Type().EmitValCall(*v, args, pos, m)	
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

func (self *LitForm) Emit(m *M) error {
	return self.val.Type().EmitVal(self.val, m)
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
