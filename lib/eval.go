package gfun

import (
	"fmt"
	"log"
	"time"
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

			if fun, err = tgt.Data(); err != nil {
				return err
			}

			retPc := pc+1
			
			if pc, err = fun.(*Fun).Call(retPc, self); err != nil {
				return err
			}

			if pc == retPc {
				self.env.outer.Regs[0] = self.env.Regs[0]
				self.env = self.env.outer
			}
			
		case BENCH:
			if (self.debug) {
				log.Printf("BENCH %v\n", op.BenchReps())
			}

			reps, err := self.env.Regs[op.BenchReps()].Data()

			if err != nil {
				return err
			}

			start := time.Now()
			
			for i := 0; i < reps.(int); i++ {
				if err := self.Eval(pc+1); err != nil {
					return err
				}
			}
			
			self.env.Regs[0].Init(&self.IntType, int(time.Since(start).Milliseconds()))
			pc = op.BenchEndPc()
			
		case BRANCH:
			if (self.debug) {
				log.Printf("BRANCH %v %v %v\n", op.BranchCond(), op.BranchTruePc(), op.BranchFalsePc())
			}
			
			if cond := self.env.Regs[op.BranchCond()]; cond.Type().BoolVal(cond) {
				pc = op.BranchTruePc()
			} else {
				pc = op.BranchFalsePc()
			}

		case DEC:
			if (self.debug) {
				log.Printf("DEC %v\n", op.DecTarget())
			}
			
			t := self.env.Regs[op.DecTarget()]
			v, err := t.Data()

			if err != nil {
				return err
			}
			
			t.Init(&self.IntType, v.(int)-op.DecDelta())
			self.env.Regs[0] = t
			pc++
		
		case ENV_POP:
			if (self.debug) {
				log.Printf("ENV_POP\n")
			}

			self.env = self.env.outer
			pc++

		case ENV_PUSH:
			if (self.debug) {
				log.Printf("ENV_PUSH\n")
			}

			self.env = new(Env).Init(self.env)
			pc++

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
			self.env.outer.Regs[0] = self.env.Regs[0]
			self.env = self.env.outer
			
		default:
			log.Fatalf("Unknown op code at pc %v: %v (%v)", pc, op.Code(), op)
		}
	}

	return nil
}
