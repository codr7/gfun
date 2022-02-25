package gfun

type IntType struct {
	BasicType
}

func (self *IntType) Name() *Sym {
	return self.m.Sym("Int")
}
