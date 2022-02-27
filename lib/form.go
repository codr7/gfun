package gfun

import (
	//"fmt"
)

type Form interface {
	Emit(in []Form, m *M) ([]Form, error)
}

type BasicForm struct {
	pos Pos
}

func (self *BasicForm) Init(pos Pos) {
	self.pos = pos
}

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

func (self *IdForm) Emit(in []Form, m *M) ([]Form, error) {
	v, err := m.env.GetVal(self.id)

	if err != nil {
		return nil, err
	}

	return v.Type().EmitVal(in, *v)
}

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

func (self *LitForm) Emit(in []Form, m *M) ([]Form, error) {
	return self.val.Type().EmitVal(in, self.val)
}
