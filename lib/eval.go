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

			fun := tgt.Data().(*Fun)
			retPc := pc+1
			var err error

			if pc, err = fun.Call(retPc, self); err != nil {
				return err
			}

			if pc == retPc {
				self.EndEnv()
			}

		case CALLI1:
			retPc := pc+1
			var err error
			
			if pc, err = op.CallI1Target().Call(retPc, self); err != nil {
				return err
			}

			if pc == retPc {
				self.EndEnv()
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
				self.EndEnv()
			}

		case COPY:
			self.Env().Regs[op.Reg1()] = self.Env().Regs[op.Reg2()]
			pc++
			
		case BENCH:
			reps := self.Env().Regs[op.BenchReps()].Data()
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
			env := self.Env()
			t := &env.Regs[op.DecTarget()]
			t.Init(&self.IntType, t.Data().(int)-op.DecDelta())
			env.Regs[0] = *t
			pc++

		case EQ:
			env := self.Env()
			l := env.Regs[op.Reg1()]
			r := env.Regs[op.Reg2()]
			env.Regs[0].Init(&self.BoolType, l.Type().EqVal(l, r))
			pc++
			
		case ENV_BEG:
			self.BeginEnv(self.Env())
			pc++

		case ENV_END:
			self.EndEnv()
			pc++

		case FUN:
			f := self.Env().Regs[op.Reg1()].Data().(*Fun)

			if err := f.CaptureEnv(self); err != nil {
				return err
			}
			
			pc = op.FunEndPc()
			
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
			
		case LOAD_NIL:
			self.Env().Regs[op.Reg1()].Init(&self.NilType, nil)
			pc++

		case LOAD_TYPE:
			t := self.types[op.LoadTypeId()]
			self.Env().Regs[op.Reg1()].Init(&self.MetaType, t)
			pc++
				
		case NOP:
			pc++

		case REC:
			prev := self.EndEnv()
			env := self.Env()

			for i := 1; i < ArgCount+1; i++ {
				env.Regs[i] = prev.Regs[i]
			}

			pc = self.Frame().startPc

		case RET:
			f := self.EndFrame()
			pc = f.retPc
			self.EndEnv()

		case TEST:
			env := self.Env()
			exp := env.Regs[op.Reg1()]
			
			if err := self.Eval(pc+1); err != nil {
				return err
			}

			if res := env.Regs[0]; !exp.Type().EqVal(exp, res) {
				return fmt.Errorf("Test failed: %v/%v", res, exp)
			}

			fmt.Printf(".")
			pc = op.TestEndPc()

		default:
			log.Fatalf("Unknown op code at pc %v: %v (%v)", pc, op.OpCode(), op)
		}
	}

	return nil
}
