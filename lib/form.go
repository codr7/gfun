package gfun

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

func (self *IdForm) Init(pos Pos, id *Sym) {
	self.BasicForm.Init(pos)
	self.id = id
}

func (self *IdForm) Emit(in []Form, m *M) ([]Form, error) {
	v, err := m.env.GetVal(self.id)

	if err != nil {
		return nil, err
	}
	
	return v.Type().EmitVal(in, *v)
}
