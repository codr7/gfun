package gfun

type Op uint64

const (
	OpCodeBits = 5
	OpRegBits = 8
	OpReg2Bits = OpCodeBits + OpRegBits
)

func (self Op) Code() int {
	return int(self & 0x0000000F)
}

func (self Op) Reg1() Reg {
	return Reg(self & 0x00000FF0 >> OpCodeBits)
}

func (self Op) Reg2() Reg {
	return Reg(self & 0x000FF000 >> OpReg2Bits)
}

func (self *M) Emit(op Op) *Op {
	pc := self.emitPc
	self.emitPc++
	self.ops[pc] = op
	return &self.ops[pc]
}

const (
	STOP = iota
	
	BRANCH_REG
	DEC
	INC
)

func (self *M) EmitBranchReg(condReg Reg) *Op {
	return self.Emit(Op(BRANCH_REG + (condReg << OpCodeBits)))
}

func (self *M) EmitDec(reg1 Reg, reg2 Reg) *Op {
	return self.Emit(Op(DEC + (reg1 << OpCodeBits) + (reg2 << OpReg2Bits)))
}

func (self *M) EmitInc(reg1 Reg, reg2 Reg) *Op {
	return self.Emit(Op(INC + (reg1 << OpCodeBits) + (reg2 << OpReg2Bits)))
}

func (self *M) EmitStop() {
	self.Emit(STOP)
}

