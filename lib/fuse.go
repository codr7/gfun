package gfun

import (
	"log"
)

func (self *M) FuseCircularCopies(startPc PC) int {
	res := 0
	
	for i := startPc; i < self.emitPc; i++ {
		op1 := self.ops[i]

		if op1.OpCode() != COPY {
			continue
		}
		
		for j := i+1; j < self.emitPc; j++ {
			op2 := &self.ops[j]

			if op2.OpCode() == COPY && op2.Reg1() == op1.Reg2() && op2.Reg2() == op1.Reg1() {
					op2.InitNop()
					log.Printf("Fused circular copy at %v: %v/%v", j, op1.Reg1(), op1.Reg2())
					res++
			} else if op2.WritesReg(op1.Reg1()) || op2.WritesReg(op1.Reg2()) {
				break
			}
		}
	}

	return res
}

func (self *M) FuseEntry(startPc PC) PC {
	done := false

	for !done && startPc < self.emitPc {
		op := self.ops[startPc]
		
		switch op.OpCode() {
		case GOTO:
			log.Printf("Fused entry to %v", op.GotoPc())
			startPc = op.GotoPc()
		case NOP:
			startPc++
			log.Printf("Fused entry to %v", startPc)
		default:
			done = true
		}
	}

	return startPc
}

func (self *M) FuseExit(startPc PC) int {
	res := 0
	done := false

	for i := self.emitPc-1; !done && i >= startPc; i-- {
		op := &self.ops[i]
		
		switch op.OpCode() {
		case GOTO, LOAD_NIL, NOP:
			log.Printf("Fused exit to %v", i)
			op.InitRet()
			res++
		case RET:
			break
		default:
			done = true
		}
	}

	return res
}

func (self *M) FuseNops(startPc PC) int {
	res := 0
	
	for i := startPc; i < self.emitPc; i++ {
		op1 := &self.ops[i]

		if op1.OpCode() != NOP {
			continue
		}
		
		for j := i+1; j < self.emitPc; j++ {
			op2 := &self.ops[j]

			if op2.OpCode() == NOP {
				op1.InitGoto(j+1)
				log.Printf("Fused nop at %v", j)
				res++
			} else if op2.OpCode() == GOTO {
				op1.InitGoto(op2.GotoPc())
				log.Printf("Fused nop at %v", j)
				res++
				break
			} else {
				break
			}
		}
	}

	return res
}

func (self *M) FuseUnusedLoads(startPc PC) int {
	res := 0
	
	for i := startPc; i < self.emitPc; i++ {
		op1 := &self.ops[i]
		reg := Reg(-1)
		
		switch op1.OpCode() {
		case COPY:
			reg = op1.Reg1()
		case LOAD_BOOL, LOAD_FUN, LOAD_INT1, LOAD_INT2, LOAD_MACRO, LOAD_TYPE:
			reg = op1.LoadTarget()
		}

		if reg != -1 {
			used := false

			if i == self.emitPc-1 {
				used = true
			}
			
			for j := i+1; j < self.emitPc; j++ {
				op2 := self.ops[j]
				
				if op2.WritesReg(reg) {
					break
				}
				
				if op2.ReadsReg(reg) || (reg == 0 && j == self.emitPc-1) || op2.OpCode() == ENV_PUSH {
					used = true
					break
				}
			}

			if !used {
				log.Printf("Fused unused load at %v: %v", i, reg) 

				switch op1.OpCode() {
				case LOAD_FUN, LOAD_INT2, LOAD_MACRO:
					self.ops[i+1].InitNop()
					i++
				}

				op1.InitNop()
				res++
			}
		}
	}

	return res
}

func (self *M) Fuse(startPc PC) PC {
	for {		
		if self.FuseCircularCopies(startPc) == 0 &&
			self.FuseUnusedLoads(startPc) == 0 &&
			self.FuseNops(startPc) == 0 &&
			self.FuseExit(startPc) == 0 {
			break
		}

		startPc = self.FuseEntry(startPc)
	}

	return startPc
}
