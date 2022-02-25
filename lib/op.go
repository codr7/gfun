package gfun

type Op uint64

func (self Op) Code() int {
	return int(self & 0x000F)
}

func (self Op) Reg() int {
	return int(self >> 4)
}

func (self *M) Emit(op Op) PC {
	pc := self.emitPc
	self.ops[pc] = op
	self.emitPc++
	return pc
}

const (
	STOP = 0
	INC = 1
)

func (self *M) EmitStop() {
	self.Emit(STOP)
}

func (self *M) EmitInc(reg int) {
	self.Emit(Op(INC + (reg << 4)))
}
