package gfun

type Op uint64

func (self Op) Code() int {
	return int(self & 0x0000000F)
}

func (self Op) Reg1() int {
	return int(self & 0x00000FF0 >> 4)
}

func (self Op) Reg2() int {
	return int(self & 0x000FF000 >> 12)
}

func (self *M) Emit(op Op) PC {
	pc := self.emitPc
	self.ops[pc] = op
	self.emitPc++
	return pc
}

const (
	STOP = 0

	ADD = 1
	DEC = 2
	INC = 3
	SUB = 4
)

func (self *M) EmitStop() {
	self.Emit(STOP)
}

func (self *M) EmitInc(reg int) {
	self.Emit(Op(INC + (reg << 4)))
}

func (self *M) EmitAdd(reg1 int, reg2 int) {
	self.Emit(Op(ADD + (reg1 << 4) + (reg2 << 12)))
}
