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

func (self *M) Emit(op Op) *Op {
	pc := self.emitPc
	self.emitPc++
	self.ops[pc] = op
	return &self.ops[pc]
}

const (
	DEC = 1
	INC = 2

	STOP = 0
)

func (self *M) EmitDec(reg1 int, reg2 int) *Op {
	return self.Emit(Op(DEC + (reg1 << 4) + (reg2 << 12)))
}

func (self *M) EmitInc(reg1 int, reg2 int) *Op {
	return self.Emit(Op(INC + (reg1 << 4) + (reg2 << 12)))
}

func (self *M) EmitStop() {
	self.Emit(STOP)
}

