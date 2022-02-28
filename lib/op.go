package gfun

import (
	"unsafe"
)

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
	LOAD_BOOL
	LOAD_FUN
	LOAD_INT1
	LOAD_INT2
	LOAD_MACRO
	MOVE
	NOP
	RET
)

func (self *M) EmitStop() {
	self.Emit(STOP)
}

const (
	OpBranchTruePcBits = OpReg2Bits
	OpBranchFalsePcBits = OpBranchTruePcBits + OpPcBits
)

func (self Op) BranchCond() Reg {
	return Reg((self >> OpCodeBits) & 1 << OpRegBits)
}

func (self Op) BranchTruePc() PC {
	return PC((self >> (OpCodeBits + OpRegBits)) & 1 << OpPcBits)
}

func (self Op) BranchFalsePc() PC {
	return PC((self >> (OpCodeBits + OpRegBits + OpPcBits)) & 1 << OpPcBits)
}

func (self *M) EmitBranch(cond Reg, truePc, falsePc PC) *Op {
	return self.Emit(Op(BRANCH + Op(cond << OpCodeBits) + Op(truePc << OpBranchTruePcBits) + Op(OpBranchFalsePcBits)))
}

/* Call */

func (self Op) CallTarget() Reg {
	return self.Reg1()
}

func (self Op) CallFlags() CallFlags {
	var flags CallFlags
	flags.Memo = (self >> OpReg2Bits) & 1 == 1
	flags.Tail = (self >> (OpReg2Bits+1)) & 1 == 1
	return flags
}

func (self *M) EmitCall(target Reg, flags CallFlags) *Op {
	op := Op(CALL + (target << OpCodeBits))

	if flags.Memo {
		op += 1 << OpReg2Bits
	}

	if flags.Tail {
		op += 1 << (OpReg2Bits+1)
	}

	return self.Emit(op)
}

/* Goto */

func (self Op) GotoPc() PC {
	return PC(self >> OpCodeBits)
}

func (self *M) EmitGoto(pc PC) *Op {
	return self.Emit(Op(GOTO + (pc << OpCodeBits)))
}

/* LoadBool */

const (
	OpLoadBoolValBits = OpReg2Bits
)

func (self Op) LoadBoolVal() bool {
	if v := self >> OpLoadBoolValBits; v == 1 {
		return true
	}

	return false
}

func (self *M) EmitLoadBool(dst Reg, val bool) *Op {
	v := 0

	if val {
		v++
	}
	
	return self.Emit(Op(LOAD_BOOL + Op(dst << OpCodeBits) + Op(v << OpLoadBoolValBits)))
}

func (self *M) EmitLoadFun(dst Reg, val *Fun) *Op {
	op := self.Emit(Op(LOAD_FUN + Op(dst << OpCodeBits)))
	self.Emit(Op(uintptr(unsafe.Pointer(val))))
	return op
}

/* LoadInt */

const (
	OpLoadInt1Max = 1 << (OpBits - OpRegBits - OpCodeBits - 1)
	OpLoadInt1ValBits = OpReg2Bits
)

func (self Op) LoadInt1Val() int {
	return int(self >> OpLoadInt1ValBits)
}

 func (self *M) EmitLoadInt(dst Reg, val int) *Op {
	if val > OpLoadInt1Max-1 || val < -OpLoadInt1Max {
		op := self.Emit(Op(LOAD_INT2 + (dst << OpCodeBits)))
		self.Emit(Op(val))
		return op
	}

	return self.Emit(Op(LOAD_INT1 + Op(dst << OpCodeBits) + Op(val << OpLoadInt1ValBits)))
}

func (self *M) EmitLoadMacro(dst Reg, val *Macro) *Op {
	op := self.Emit(Op(LOAD_MACRO + Op(dst << OpCodeBits)))
	self.Emit(Op(uintptr(unsafe.Pointer(val))))
	return op
}

func (self *M) EmitMove(dst Reg, src int) *Op {
	return self.Emit(Op(MOVE + Op(dst << OpCodeBits) + Op(src << OpReg2Bits)))
}

func (self *M) EmitNop() {
	self.Emit(NOP)
}

func (self *M) EmitRet() {
	self.Emit(RET)
}
