package gfun

type CallFlags struct {
	Drop, Memo, Tail bool
}

type FunBody = func(*Fun, CallFlags, PC) (PC, error)

type Fun struct {
	m *M
	name *Sym
	body FunBody
}

func (self *Fun) Init(m *M, name *Sym, body FunBody) {
	self.m = m
	self.name = name
	self.body = body
}

func (self *Fun) Call(flags CallFlags, ret PC) (PC, error) {
	return self.body(self, flags, ret)
}
