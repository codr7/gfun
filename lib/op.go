package gfun

import (
	"fmt"
	"io"
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

func (self Op) OpCode() int {
	return int(self & ((1 << OpCodeBits) - 1))
}

func (self Op) Reg1() Reg {
	return Reg((self >> OpCodeBits) & ((1 << OpRegBits) - 1))
}

func (self Op) Reg2() Reg {
	return Reg((self >> OpReg2Bit) & ((1 << OpRegBits) - 1))
}

func (self Op) ReadsReg(reg Reg) bool {
	switch self.OpCode() {
	case BENCH:
		if self.BenchReps() == reg {
			return true
		}
	case BRANCH:
		if self.BranchCond() == reg {
			return true
		}
	case CALL:
		if (reg > 0 && reg < ArgCount+1) || reg == self.CallTarget(){
			return true
		}
	case CALLI1, CALLI2, REC:
		if (reg > 0 && reg < ArgCount+1) {
			return true
		}
	case COPY:
		if self.Reg2() == reg {
			return true
		}
	case EQ:
		if self.Reg1() == reg || self.Reg2() == reg {
			return true
		}
	case RET:
		if reg == 0 {
			return true
		}
	case TEST:
		if self.Reg1() == reg {
			return true
		}
	default:
		break
	}

	return false
}

func (self Op) WritesReg(reg Reg) bool {
	switch self.OpCode() {
	case CALL, CALLI1, CALLI2:
		if reg == 0 {
			return true
		}
	case COPY:
		if self.Reg1() == reg {
			return true
		}
	case DEC:
		if self.DecTarget() == reg || reg == 0 {
			return true
		}
	case LOAD_BOOL, LOAD_FUN, LOAD_INT1, LOAD_INT2, LOAD_MACRO, LOAD_NIL, LOAD_TYPE:
		if self.LoadTarget() == reg {
			return true
		}
	default:
		break
	}

	return false
}

func (self Op) Dump(pc PC, m *M, out io.Writer) PC {
		switch self.OpCode() {
		case STOP:
			fmt.Fprintf(out, "STOP")
		case CALL:
			fmt.Fprintf(out, "CALL %v", self.CallTarget())
		case CALLI1:
			fmt.Fprintf(out, "CALLI1 %v", self.CallI1Target())
		case CALLI2:
			d := unsafe.Pointer(uintptr(m.ops[pc+1]))
			f := (*Fun)(d)
			fmt.Fprintf(out, "CALLI2 %v", f)
			return pc+1
		case COPY:
			fmt.Fprintf(out, "COPY %v %v", self.Reg1(), self.Reg2())
		case BENCH:
			fmt.Fprintf(out, "BENCH %v", self.BenchReps())
		case BRANCH:
			fmt.Fprintf(out, "BRANCH %v %v %v", self.BranchCond(), self.BranchTruePc(), self.BranchFalsePc())
		case DEC:
			fmt.Fprintf(out, "DEC %v %v", self.DecTarget(), self.DecDelta())
		case ENV_BEG:
			fmt.Fprintf(out, "ENV_BEG")
		case ENV_END:
			fmt.Fprintf(out, "ENV_END")
		case EQ:
			fmt.Fprintf(out, "EQ %v %v", self.Reg1(), self.Reg2())
		case FUN:
			fmt.Fprintf(out, "FUN %v %v", self.Reg1(), self.FunEndPc())
		case GOTO:
			fmt.Fprintf(out, "GOTO %v", self.GotoPc())
		case LOAD_BOOL:
			fmt.Fprintf(out, "LOAD_BOOL %v %v", self.Reg1(), self.LoadBoolVal())
		case LOAD_FUN:
			d := unsafe.Pointer(uintptr(m.ops[pc+1]))
			f := (*Fun)(d)
			fmt.Fprintf(out, "LOAD_FUN %v %v", self.Reg1(), f)
			return pc+1
		case LOAD_INT1:
			fmt.Fprintf(out, "LOAD_INT1 %v %v", self.Reg1(), self.LoadInt1Val())
		case LOAD_INT2:
			val := int(m.ops[pc+1])
			fmt.Fprintf(out, "LOAD_INT2 %v %v", self.Reg1(), val)
			return pc+1
		case LOAD_MACRO:
			d := unsafe.Pointer(uintptr(m.ops[pc+1]))
			m := (*Macro)(d)
			fmt.Fprintf(out, "LOAD_MACRO %v %v", self.Reg1(), m)
			return pc+1
		case LOAD_NIL:
			fmt.Fprintf(out, "LOAD_NIL %v", self.Reg1())
		case LOAD_TYPE:
			t := m.types[self.LoadTypeId()]
			fmt.Fprintf(out, "LOAD_TYPE %v %v", self.Reg1(), t)
		case NOP:
			fmt.Fprintf(out, "NOP")
		case REC:
			fmt.Fprintf(out, "REC")
		case RET:
			fmt.Fprintf(out, "RET")
		case TEST:
			fmt.Fprintf(out, "TEST %v %v", self.Reg1(), self.TestEndPc())
		default:
			break
		}

	return pc
}

func (self *M) DumpOps(startPc PC, out io.Writer) {
	for i := startPc; i < self.emitPc; i++ {
		fmt.Fprintf(out, "%v ", i)
		i = self.ops[i].Dump(i, self, out)
		fmt.Fprintf(out, "\n")
	}
}

func (self *M) Emit(op Op) *Op {
	pc := self.emitPc
	self.emitPc++
	self.ops[pc] = op
	return &self.ops[pc]
}

const (
	STOP = iota

	BENCH
	BRANCH
	CALL
	CALLI1
	CALLI2
	COPY
	DEC
	ENV_BEG
	ENV_END
	EQ
	FUN
	GOTO
	LOAD_BOOL
	LOAD_FUN
	LOAD_INT1
	LOAD_INT2
	LOAD_MACRO
	LOAD_NIL
	LOAD_TYPE
	NOP
	REC
	RET
	TEST
)

func (self *M) EmitStop() {
	self.Emit(STOP)
}

/* Bench */

const (
	OpBenchRepsBit = OpCodeBits
	OpBenchEndPcBit = OpBenchRepsBit+OpRegBits
)

func (self Op) BenchReps() Reg {
	return Reg((self >> OpBenchRepsBit) & ((1 << OpRegBits) - 1))
}

func (self Op) BenchEndPc() PC {
	return PC((self >> OpBenchEndPcBit) & ((1 << OpPcBits) - 1))
}

func (self *Op) InitBench(reps Reg, endPc PC) *Op {
	*self = Op(BENCH +
		Op(reps << OpBenchRepsBit) +
		Op(endPc << OpBenchEndPcBit))
	
	return self
}

func (self *M) EmitBench(reps Reg, endPc PC) *Op {
	return self.Emit(0).InitBench(reps, endPc)
}

/* Branch */

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

/* CallI */

const (
	OpCallI1TargetBits = OpBits - OpCodeBits
	OpCallI1TargetMax = 1 << OpCallI1TargetBits
)

func (self *M) EmitCallI(target *Fun) *Op {
	if tp := uintptr(unsafe.Pointer(target)); tp < OpCallI1TargetMax {
		return self.EmitCallI1(target)
	}
	
	return self.EmitCallI2(target)
}

func (self Op) CallI1Target() *Fun {
	return (*Fun)(unsafe.Pointer(uintptr((self >> OpCodeBits))))
}

func (self *M) EmitCallI1(target *Fun) *Op {
	tp := uintptr(unsafe.Pointer(target))
	return self.Emit(Op(CALLI1) + Op(tp << OpCodeBits))
}

func (self *M) EmitCallI2(target *Fun) *Op {
	op := self.Emit(Op(CALLI2))
	self.Emit(Op(uintptr(unsafe.Pointer(target))))
	return op
}

func (self *Op) InitCopy (dst Reg, src Reg) *Op {
	*self = Op(COPY + Op(dst << OpCodeBits) + Op(src << OpReg2Bit))
	return self
}

func (self *M) EmitCopy(dst Reg, src Reg) *Op {
	return self.Emit(0).InitCopy(dst, src)
}

/* Dec */

const (
	OpDecTargetBit = OpCodeBits
	OpDecDeltaBit = OpDecTargetBit + OpRegBits
	OpDecDeltaBits = OpBits - OpDecDeltaBit
)

func (self Op) DecTarget() Reg {
	return Reg((self >> OpDecTargetBit) & ((1 << OpRegBits) - 1))
}

func (self Op) DecDelta() int {
	return int((self >> OpDecDeltaBit) & ((1 << OpDecDeltaBits) - 1))
}

func (self *M) EmitDec(target Reg, delta int) *Op {
	return self.Emit(Op(DEC + Op(target << OpDecTargetBit) + Op(delta << OpDecDeltaBit)))
}

func (self *M) EmitEnvBeg() {
	self.Emit(ENV_BEG)
}

func (self *M) EmitEnvEnd() {
	self.Emit(ENV_END)
}

/* Eq */

func (self *Op) InitEq (l, r Reg) *Op {
	*self = Op(EQ + Op(l << OpCodeBits) + Op(r << OpReg2Bit))
	return self
}

func (self *M) EmitEq(l, r Reg) *Op {
	return self.Emit(0).InitEq(l, r)
}

/* Fun */

func (self Op) FunEndPc() PC {
	return PC((self >> OpReg2Bit) & ((1 << OpPcBits) - 1))
}

func (self *Op) InitFun (dst Reg, endPc PC) *Op {
	*self = Op(FUN + Op(dst << OpCodeBits) + Op(endPc << OpReg2Bit))
	return self
}

func (self *M) EmitFun(dst Reg, endPc PC) *Op {
	return self.Emit(0).InitFun(dst, endPc)
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

func (self Op) LoadTarget() Reg {
	return self.Reg1()
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

func (self *M) EmitLoadNil(dst Reg) *Op {
	return self.Emit(Op(LOAD_NIL))
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

func (self *Op) InitNop() *Op {
	*self = NOP
	return self
}

func (self *M) EmitNop() *Op {
	return self.Emit(0).InitNop()
}

func (self *Op) InitRec() *Op {
	*self = REC
	return self
}

func (self *M) EmitRec() * Op {
	return self.Emit(0).InitRec()
}

func (self *Op) InitRet() *Op {
	*self = RET
	return self
}

func (self *M) EmitRet() *Op {
	return self.Emit(0).InitRet()
}

func (self Op) TestEndPc() PC {
	return PC((self >> OpReg2Bit) & ((1 << OpPcBits) - 1))
}

func (self *Op) InitTest(exp Reg, endPc PC) *Op {
	*self = Op(TEST + Op(exp << OpCodeBits) + Op(endPc << OpReg2Bit))
	return self
}

func (self *M) EmitTest(exp Reg, endPc PC) *Op {
	return self.Emit(0).InitTest(exp, endPc)
}
