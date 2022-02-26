package gfun

type Form interface {
	Emit(m *M) error
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

func (self *IdForm) Init(pos Pos, id *Sym) {
	self.BasicForm.Init(pos)
	self.id = id
}

func (self *IdForm) Emit(m *M) error {
	v, err := m.env.GetVal(self.id)

	if err != nil {
		return err
	}
	
	return v.Type().EmitVal(*v)
}
