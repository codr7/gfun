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
					res++
					log.Printf("Fused circular copy at %v: %v/%v", j, op1.Reg1(), op1.Reg2())
			} else if op2.WritesReg(op1.Reg1()) || op2.WritesReg(op1.Reg2()) {
				break
			}
		}
	}

	return res
}

func (self *M) FuseNops(startPc PC) int {
	return 0
}

func (self *M) FuseUnusedRegs(startPc PC) int {
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
				log.Printf("Fused unused reg at %v: %v", i, reg) 

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

func (self *M) Fuse(startPc PC) {
	for self.FuseCircularCopies(startPc) > 0 ||
		self.FuseUnusedRegs(startPc) > 0 ||
		self.FuseNops(startPc) > 0 {
	}
}
