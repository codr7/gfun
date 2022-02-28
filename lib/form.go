package gfun

import (
	"fmt"
)

type Form interface {
	Emit(*M) error
	EmitCall(CallFlags, *M) error
}

type BasicForm struct {
	pos Pos
}

func (self *BasicForm) Init(pos Pos) {
	self.pos = pos
}

func (self *BasicForm) EmitCall(flags CallFlags, m *M) error {
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
	var flags CallFlags
	var err error

	for i, a := range self.args {
		if err = a.Emit(m); err != nil {
			return err
		}

		m.EmitMove(Reg(i+1), 0)
	}
	
	if  err = self.target.EmitCall(flags, m); err != nil {
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

func (self *IdForm) EmitCall(flags CallFlags, m *M) error {
	v, err := m.env.GetVal(self.id)

	if err != nil {
		return err
	}

	return v.Type().EmitValCall(*v, flags, m)	
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

func (self *LitForm) EmitCall(flags CallFlags, m *M) error {
	return self.val.Type().EmitValCall(self.val, flags, m)
}
