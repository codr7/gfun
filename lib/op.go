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
	INC
	LOAD_INT1
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

func (self *M) EmitInc(dst Reg, src Reg) *Op {
	return self.Emit(Op(INC + (dst << OpCodeBits) + (src << OpReg2Bits)))
}

const (
	OpLoadInt1Max = 1 << (64 - OpRegBits - OpCodeBits - 1)
	OpLoadInt1ValBits = OpReg2Bits
)

func (self Op) LoadInt1Val() int {
	return int(self >> OpLoadInt1ValBits)
}

func (self *M) EmitLoadInt(dst Reg, val int) *Op {
	if val > OpLoadInt1Max-1 || val < -OpLoadInt1Max {
		self.Emit(Op(LOAD_INT2 + (dst << OpCodeBits)))
		return self.Emit(Op(val))
	}

	return self.Emit(Op(LOAD_INT1 + Op(dst << OpCodeBits) + Op(val << OpLoadInt1ValBits)))
}

func (self *M) EmitStop() {
	self.Emit(STOP)
}

