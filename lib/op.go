package gfun

type Op uint64

const (
	OpCodeBits = 6
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
	
	BRANCH
	CALL
	DEC
	INC
	LOAD_INT2
)

func (self *M) EmitBranch(cond Reg) *Op {
	return self.Emit(Op(BRANCH + (cond << OpCodeBits)))
}

func (self *M) EmitCall(target Reg, flags CallFlags) *Op {
	op := Op(CALL + (target << OpCodeBits))

	if flags.Drop {
		op += 1 << OpReg2Bits
	}

	if flags.Memo {
		op += 1 << (OpReg2Bits+1)
	}

	if flags.Tail {
		op += 1 << (OpReg2Bits+2)
	}

	return self.Emit(op)
}

func (self *M) EmitDec(dst Reg, src Reg) *Op {
	return self.Emit(Op(DEC + (dst << OpCodeBits) + (src << OpReg2Bits)))
}

func (self *M) EmitInc(dst Reg, src Reg) *Op {
	return self.Emit(Op(INC + (dst << OpCodeBits) + (src << OpReg2Bits)))
}

func (self *M) EmitLoadInt(dst Reg, val int) *Op {
	self.Emit(Op(LOAD_INT2 + (dst << OpCodeBits)))
	return self.Emit(Op(val))
}

func (self *M) EmitStop() {
	self.Emit(STOP)
}

