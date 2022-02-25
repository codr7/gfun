package gfun

import (
	"fmt"
	"log"
)

func (self *M) Eval(pc PC) error {
	env := &self.RootEnv
	
	for {
		op := self.ops[pc]
		
		switch op.Code() {
		case CALL:
			fmt.Printf("CALL %v\n", op.Reg1())
			tgt := env.Regs[op.Reg1()]
			
			if tgt.Type() != &self.FunType {
				return fmt.Errorf("Not callable: %v", tgt)
			}

			var flags CallFlags
			flags.Drop = (op >> OpReg2Bits) & 01 == 1
			flags.Drop = (op >> (OpReg2Bits+1)) & 01 == 1
			flags.Drop = (op >> (OpReg2Bits+2)) & 01 == 1

			var err error
			var fun interface{}
			
			fun, err = tgt.Data()

			if err != nil {
				return err
			}

			var ret PC
			ret, err = fun.(*Fun).Call(flags, pc+1)

			if err != nil {
				return err
			}

			pc = ret
			
		case BRANCH:
			fmt.Printf("BRANCH %v\n", op.Reg1())
			cond := env.Regs[op.Reg1()]
			res, err := cond.Type().BoolVal(cond);

			if err != nil {
				return err
			}

			if res {
				pc++
			} else {
				pc += 2
			}

		case INC:
			fmt.Printf("INC %v %v\n", op.Reg1(), op.Reg2())
			l := &env.Regs[op.Reg1()]
			var lv interface{}
			var err error
			
			if lv, err = l.Data(); err != nil {
				return err
			}
			
			r := env.Regs[op.Reg2()]
			var rv interface{}
			
			if rv, err = r.Data(); err != nil {
				return err
			}
			
			l.Init(l.Type(), lv.(int)+rv.(int))
			pc++

		case LOAD_INT2:
			val := int(self.ops[pc+1])
			fmt.Printf("LOAD_INT2 %v %v\n", op.Reg1(), val)
			env.Regs[op.Reg1()].Init(&self.IntType, val)
			pc += 2
			
		case STOP:
			fmt.Printf("STOP\n")
			return nil

		default:
			log.Fatalf("Unknown op code at pc %v: %v (%v)", pc, op.Code(), op)
		}
	}

	return nil
}
