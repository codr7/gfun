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

func (self *Fun) Emit(body Form) (PC, error) {
	startPc := self.m.emitPc

	if err := body.Emit(self.m); err != nil {
		return -1, err
	}

	self.m.EmitRet()
	
	self.body = func(fun *Fun, flags CallFlags, ret PC) (PC, error) {
		self.m.Call(fun, flags, ret)
		return startPc, nil
	}

	return startPc, nil
}
