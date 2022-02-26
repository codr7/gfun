package gfun

type Val struct {
	_type Type
	data interface{}
}

func NewVal(_type Type, data interface{}) Val {
	var self Val
	self.Init(_type, data)
	return self
}

func (self *Val) Init(_type Type, data interface{}) {
	self._type = _type
	self.data = data
}

func (self *Val) Type() Type {
	return self._type
}

func (self *Val) Data() (interface{}, error) {
	return self._type.GetVal(self.data)
}
