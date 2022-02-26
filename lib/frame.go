package gfun

type Frame struct {
	Env Env
	
	outer *Frame
	fun *Fun
	callFlags CallFlags
	ret PC
}

func (self *Frame) Init(m *M, fun *Fun, callFlags CallFlags, ret PC) *Frame {
	self.outer = m.frame
	self.Env.Init(m.env)
	self.fun = fun
	self.callFlags = callFlags
	self.ret = ret
	return self
}

func (self *M) Call(fun *Fun, callFlags CallFlags, ret PC) *Frame {
	self.frame = new(Frame).Init(self, fun, callFlags, ret)
	return self.frame
}

func (self *M) Ret() *Frame {
	f := self.frame
	self.frame = self.frame.outer
	return f
}
