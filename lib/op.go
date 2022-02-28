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
	OpReg2Bit = OpCodeBits + OpRegBits
	OpTypeIdBits = 10
)

func (self Op) Code() int {
	return int(self & ((1 << OpCodeBits)-1))
}

func (self Op) Reg1() Reg {
	return Reg(self & ((1 << OpReg2Bit)-1) >> OpCodeBits)
}

func (self Op) Reg2() Reg {
	return Reg(self & ((1 << (OpReg2Bit+OpRegBits))-1) >> OpReg2Bit)
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
	LOAD_TYPE
	COPY
	NOP
	RET
)

func (self *M) EmitStop() {
	self.Emit(STOP)
}

const (
	OpBranchTruePcBit = OpReg2Bit
	OpBranchFalsePcBit = OpBranchTruePcBit + OpPcBits
)

func (self Op) BranchCond() Reg {
	return Reg((self >> OpCodeBits) & ((1 << OpRegBits) - 1))
}

func (self Op) BranchTruePc() PC {
	return PC((self >> OpBranchTruePcBit) & ((1 << OpPcBits) - 1))
}

func (self Op) BranchFalsePc() PC {
	return PC((self >> OpBranchFalsePcBit) & ((1 << OpPcBits) - 1))
}

func (self *Op) InitBranch(cond Reg, truePc, falsePc PC) *Op {
	*self = Op(BRANCH + Op(cond << OpCodeBits) +
		Op(truePc << OpBranchTruePcBit) +
		Op(falsePc << OpBranchFalsePcBit))
	
	return self
}

func (self *M) EmitBranch(cond Reg, truePc, falsePc PC) *Op {
	return self.Emit(0).InitBranch(cond, truePc, falsePc)
}

/* Call */

func (self Op) CallTarget() Reg {
	return self.Reg1()
}

func (self *M) EmitCall(target Reg) *Op {
	return self.Emit(Op(CALL + (target << OpCodeBits)))
}

/* Goto */

func (self Op) GotoPc() PC {
	return PC(self >> OpCodeBits)
}

func (self *Op) InitGoto(pc PC) *Op {
	*self = Op(GOTO + (pc << OpCodeBits))
	return self
}

func (self *M) EmitGoto(pc PC) *Op {
	return self.Emit(0).InitGoto(pc)
}

/* LoadBool */

const (
	OpLoadBoolValBit = OpReg2Bit
)

func (self Op) LoadBoolVal() bool {
	if v := self >> OpLoadBoolValBit; v == 1 {
		return true
	}

	return false
}

func (self *M) EmitLoadBool(dst Reg, val bool) *Op {
	v := 0

	if val {
		v++
	}
	
	return self.Emit(Op(LOAD_BOOL + Op(dst << OpCodeBits) + Op(v << OpLoadBoolValBit)))
}

func (self *M) EmitLoadFun(dst Reg, val *Fun) *Op {
	op := self.Emit(Op(LOAD_FUN + Op(dst << OpCodeBits)))
	self.Emit(Op(uintptr(unsafe.Pointer(val))))
	return op
}

/* LoadInt */

const (
	OpLoadInt1Max = 1 << (OpBits - OpRegBits - OpCodeBits - 1)
	OpLoadInt1ValBit = OpReg2Bit
)

func (self Op) LoadInt1Val() int {
	v := int(self >> OpLoadInt1ValBit)
	
	if v >= OpLoadInt1Max {
		return v - (OpLoadInt1Max << 1)
	}
	
	return v
}

func (self *M) EmitLoadInt(dst Reg, val int) *Op {
	if val > OpLoadInt1Max-1 || val < -OpLoadInt1Max {
		op := self.Emit(Op(LOAD_INT2 + (dst << OpCodeBits)))
		self.Emit(Op(val))
		return op
	}

	return self.Emit(Op(LOAD_INT1 + Op(dst << OpCodeBits) + Op(val << OpLoadInt1ValBit)))
}

func (self *M) EmitLoadMacro(dst Reg, val *Macro) *Op {
	op := self.Emit(Op(LOAD_MACRO + Op(dst << OpCodeBits)))
	self.Emit(Op(uintptr(unsafe.Pointer(val))))
	return op
}

const (
	OpLoadTypeIdBit = OpReg2Bit
)

func (self Op) LoadTypeId() TypeId {
	return TypeId((self >> OpLoadTypeIdBit) & ((1 << OpTypeIdBits) - 1))
}

func (self *M) EmitLoadType(dst Reg, _type Type) *Op {
	return self.Emit(Op(LOAD_TYPE + Op(dst << OpCodeBits) + Op(_type.Id() << OpLoadTypeIdBit)))
}

func (self *M) EmitCopy(dst Reg, src Reg) *Op {
	return self.Emit(Op(COPY + Op(dst << OpCodeBits) + Op(src << OpReg2Bit)))
}

func (self *M) EmitNop() {
	self.Emit(NOP)
}

func (self *M) EmitRet() {
	self.Emit(RET)
}
