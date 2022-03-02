package gfun

type Frame struct {
	outer *Frame
	fun *Fun
	startPc, retPc PC
}

func (self *Frame) Init(outer *Frame, fun *Fun, startPc, retPc PC) *Frame {
	self.fun = fun
	self.startPc = startPc
	self.retPc = retPc
	return self
}

func (self *M) Frame() *Frame {
	if self.frameCount == 0 {
		return nil
	}
	
	return &self.frames[self.frameCount-1]
}

func (self *M) PushFrame(fun *Fun, startPc, retPc PC) *Frame {
	f := &self.frames[self.frameCount]
	f.Init(self.Frame(), fun, startPc, retPc)
	self.frameCount++
	return f
}

func (self *M) PopFrame() *Frame {
	self.frameCount--
	return &self.frames[self.frameCount]
}
