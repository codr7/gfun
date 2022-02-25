package gfun

type Op uint64

const (
	OpBits = 64
	OpCodeBits = 6
	OpPcBits = 10
	OpRegBits = 8
	OpReg2Bits = OpCodeBits + OpRegBits
)

func (self Op) Code() int {
	return int(self & ((1 << OpCodeBits)-1))
}

func (self Op) Reg1() Reg {
	return Reg(self & ((1 << OpReg2Bits)-1) >> OpCodeBits)
}

func (self Op) Reg2() Reg {
	return Reg(self & ((1 << (OpReg2Bits+OpRegBits))-1) >> OpReg2Bits)
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
	GOTO
	INC
	LOAD_INT1
	LOAD_INT2
	NOP
)

func (self *M) EmitStop() {
	self.Emit(STOP)
}

func (self *M) EmitBranch(cond Reg) *Op {
	return self.Emit(Op(BRANCH + (cond << OpCodeBits)))
}

func (self Op) CallFlags() CallFlags {
	var flags CallFlags
	flags.Drop = (self >> OpReg2Bits) & 1 == 1
	flags.Drop = (self >> (OpReg2Bits+1)) & 1 == 1
	flags.Drop = (self >> (OpReg2Bits+2)) & 1 == 1
	return flags
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

const (
	OpGotoPcBits = OpCodeBits
)

func (self Op) GotoPc() PC {
	return PC(self >> OpGotoPcBits)
}

func (self *M) EmitGoto(pc PC) *Op {
	return self.Emit(Op(GOTO + (pc << OpCodeBits)))
}

func (self *M) EmitInc(dst Reg, src Reg) *Op {
	return self.Emit(Op(INC + (dst << OpCodeBits) + (src << OpReg2Bits)))
}

const (
	OpLoadInt1Max = 1 << (OpBits - OpRegBits - OpCodeBits - 1)
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

func (self *M) EmitNop() {
	self.Emit(NOP)
}
