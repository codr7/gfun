package gfun

type Val struct {
	gtype Type
	data interface{}
}

func (self *Val) Init(gtype Type, data interface{}) {
	self.gtype = gtype
	self.data = data
}

func (self *Val) Type() Type {
	return self.gtype
}

func (self *Val) Data() (interface{}, error) {
	return self.gtype.GetVal(self.data)
}
