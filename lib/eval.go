package gfun

import (
	"fmt"
	"log"
	"os"
	"time"
	"unsafe"
)

func (self *M) Eval(pc PC) error {
	for {
		op := self.ops[pc]

		if (self.debug) {
			op.Dump(pc, self, os.Stdout)
			fmt.Println("")
		}

		switch op.OpCode() {
		case STOP:
			return nil
			
		case CALL:
			tgt := self.Env().Regs[op.CallTarget()]
			
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
				self.PopEnv()
			}

		case CALLI1:
			retPc := pc+1
			var err error
			
			if pc, err = op.CallI1Target().Call(retPc, self); err != nil {
				return err
			}

			if pc == retPc {
				self.PopEnv()
			}

		case CALLI2:
			d := unsafe.Pointer(uintptr(self.ops[pc+1]))
			f := (*Fun)(d)

			retPc := pc+1
			var err error
			
			if pc, err = f.Call(retPc, self); err != nil {
				return err
			}

			if pc == retPc {
				self.PopEnv()
			}

		case COPY:
			self.Env().Regs[op.Reg1()] = self.Env().Regs[op.Reg2()]
			pc++
			
		case BENCH:
			reps, err := self.Env().Regs[op.BenchReps()].Data()

			if err != nil {
				return err
			}

			start := time.Now()
			
			for i := 0; i < reps.(int); i++ {
				if err := self.Eval(pc+1); err != nil {
					return err
				}
			}
			
			self.Env().Regs[0].Init(&self.IntType, int(time.Since(start).Milliseconds()))
			pc = op.BenchEndPc()
			
		case BRANCH:
			if cond := self.Env().Regs[op.BranchCond()]; cond.Type().BoolVal(cond) {
				pc = op.BranchTruePc()
			} else {
				pc = op.BranchFalsePc()
			}

		case DEC:
			t := self.Env().Regs[op.DecTarget()]
			v, err := t.Data()

			if err != nil {
				return err
			}
			
			t.Init(&self.IntType, v.(int)-op.DecDelta())
			self.Env().Regs[0] = t
			pc++
		
		case ENV_POP:
			self.PopEnv()
			pc++

		case ENV_PUSH:
			self.PushEnv()
			pc++

		case GOTO:
			pc = op.GotoPc()
			
		case LOAD_BOOL:
			self.Env().Regs[op.Reg1()].Init(&self.BoolType, op.LoadBoolVal())
			pc++

		case LOAD_FUN:
			d := unsafe.Pointer(uintptr(self.ops[pc+1]))
			f := (*Fun)(d)
			self.Env().Regs[op.Reg1()].Init(&self.FunType, f)
			pc += 2

		case LOAD_INT1:
			self.Env().Regs[op.Reg1()].Init(&self.IntType, op.LoadInt1Val())
			pc++

		case LOAD_INT2:
			val := int(self.ops[pc+1])
			self.Env().Regs[op.Reg1()].Init(&self.IntType, val)
			pc += 2

		case LOAD_MACRO:
			d := unsafe.Pointer(uintptr(self.ops[pc+1]))
			m := (*Macro)(d)
			self.Env().Regs[op.Reg1()].Init(&self.MacroType, m)
			pc += 2
			
		case LOAD_TYPE:
			t := self.types[op.LoadTypeId()]
			self.Env().Regs[op.Reg1()].Init(&self.MetaType, t)
			pc++
				
		case NOP:
			pc++

		case RET:
			f := self.PopFrame()
			pc = f.ret
			self.PopEnv()
			
		default:
			log.Fatalf("Unknown op code at pc %v: %v (%v)", pc, op.OpCode(), op)
		}
	}

	return nil
}
