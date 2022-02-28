package gfun

type Frame struct {
	Env Env
	
	outer *Frame
	fun *Fun
	ret PC
}

func (self *Frame) Init(m *M, fun *Fun, ret PC) *Frame {
	self.outer = m.frame
	self.Env.Init(m.env)
	self.fun = fun
	self.ret = ret
	return self
}

func (self *M) Call(fun *Fun, ret PC) *Frame {
	f := new(Frame).Init(self, fun, ret)
	self.frame = f
	self.env = &f.Env
	return self.frame
}

func (self *M) Ret() *Frame {
	f := self.frame
	self.frame = f.outer
	self.env.outer.Regs[0] = self.env.Regs[0]
	self.env = self.env.outer
	return f
}
