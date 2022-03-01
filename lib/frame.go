package gfun

type Frame struct {
	outer *Frame
	fun *Fun
	ret PC
}

func (self *Frame) Init(outer *Frame, fun *Fun, ret PC) *Frame {
	self.fun = fun
	self.ret = ret
	return self
}

func (self *M) Frame() *Frame {
	if self.frameCount == 0 {
		return nil
	}
	
	return &self.frames[self.frameCount-1]
}

func (self *M) PushFrame(fun *Fun, ret PC) *Frame {
	f := &self.frames[self.frameCount]
	f.Init(self.Frame(), fun, ret)
	self.frameCount++
	return f
}

func (self *M) PopFrame() *Frame {
	self.frameCount--
	return &self.frames[self.frameCount]
}
