package gfun

import (
	"fmt"
	"log"
	"unsafe"
)

func (self *M) Eval(pc PC) error {
	for {
		op := self.ops[pc]
		
		switch op.Code() {
		case STOP:
			if (self.debug) {
				log.Printf("STOP\n")
			}
			
			return nil
			
		case CALL:
			if (self.debug) {
				log.Printf("CALL %v\n", op.CallTarget())
			}
			
			tgt := self.env.Regs[op.CallTarget()]
			
			if tgt.Type() != &self.FunType {
				return fmt.Errorf("Not callable: %v", tgt)
			}

			var err error
			var fun interface{}
			fun, err = tgt.Data()

			if err != nil {
				return err
			}

			pc, err = fun.(*Fun).Call(pc+1, self)

			if err != nil {
				return err
			}
			
		case BRANCH:
			if (self.debug) {
				log.Printf("BRANCH %v %v %v\n", op.BranchCond(), op.BranchTruePc(), op.BranchFalsePc())
			}
			
			if cond := self.env.Regs[op.BranchCond()]; cond.Type().BoolVal(cond) {
				pc = op.BranchTruePc()
			} else {
				pc = op.BranchFalsePc()
			}

		case GOTO:
			if (self.debug) {
				log.Printf("GOTO %v\n", op.GotoPc())
			}
			
			pc = op.GotoPc()
			
		case LOAD_BOOL:
			if (self.debug) {
				log.Printf("LOAD_BOOL %v %v\n", op.Reg1(), op.LoadBoolVal())
			}
			
			self.env.Regs[op.Reg1()].Init(&self.BoolType, op.LoadBoolVal())
			pc++

		case LOAD_FUN:
			d := unsafe.Pointer(uintptr(self.ops[pc+1]))
			f := (*Fun)(d)

			if (self.debug) {
				log.Printf("LOAD_FUN %v %v\n", op.Reg1(), f)
			}
			
			self.env.Regs[op.Reg1()].Init(&self.FunType, f)
			pc += 2

		case LOAD_INT1:
			if (self.debug) {
				log.Printf("LOAD_INT1 %v %v\n", op.Reg1(), op.LoadInt1Val())
			}
			
			self.env.Regs[op.Reg1()].Init(&self.IntType, op.LoadInt1Val())
			pc++

		case LOAD_INT2:
			val := int(self.ops[pc+1])

			if (self.debug) {
				log.Printf("LOAD_INT2 %v %v\n", op.Reg1(), val)
			}
			
			self.env.Regs[op.Reg1()].Init(&self.IntType, val)
			pc += 2

		case LOAD_MACRO:
			d := unsafe.Pointer(uintptr(self.ops[pc+1]))
			m := (*Macro)(d)

			if (self.debug) {
				log.Printf("LOAD_MACRO %v %v\n", op.Reg1(), m)
			}
			
			self.env.Regs[op.Reg1()].Init(&self.MacroType, m)
			pc += 2
			
		case LOAD_TYPE:
			t := self.types[op.LoadTypeId()]

			if (self.debug) {
				log.Printf("LOAD_TYPE %v\n", op.Reg1(), t)
			}

			self.env.Regs[op.Reg1()].Init(&self.MetaType, t)
			pc++

		case COPY:
			if (self.debug) {
				log.Printf("COPY %v %v\n", op.Reg1(), op.Reg2())
			}
			
			self.env.Regs[op.Reg1()] = self.env.Regs[op.Reg2()]
			pc++
				
		case NOP:
			if (self.debug) {
				log.Printf("NOP")
			}
			
			pc++

		case RET:
			if (self.debug) {
				log.Printf("RET")
			}
			
			f := self.Ret()
			pc = f.ret
			
		default:
			log.Fatalf("Unknown op code at pc %v: %v (%v)", pc, op.Code(), op)
		}
	}

	return nil
}
