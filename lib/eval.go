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
		case STOP:
			fmt.Printf("STOP\n")
			return nil
			
		case CALL:
			fmt.Printf("CALL %v\n", op.Reg1())
			tgt := env.Regs[op.Reg1()]
			
			if tgt.Type() != &self.FunType {
				return fmt.Errorf("Not callable: %v", tgt)
			}

			var err error
			var fun interface{}
			fun, err = tgt.Data()

			if err != nil {
				return err
			}

			var ret PC
			ret, err = fun.(*Fun).Call(op.CallFlags(), pc+1)

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

		case GOTO:
			fmt.Printf("GOTO %v\n", op.GotoPc())			
			pc = op.GotoPc()
			
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

		case LOAD_INT1:
			fmt.Printf("LOAD_INT1 %v %v\n", op.Reg1(), op.LoadInt1Val())
			env.Regs[op.Reg1()].Init(&self.IntType, op.LoadInt1Val())
			pc++

		case LOAD_INT2:
			val := int(self.ops[pc+1])
			fmt.Printf("LOAD_INT2 %v %v\n", op.Reg1(), val)
			env.Regs[op.Reg1()].Init(&self.IntType, val)
			pc += 2

		case NOP:
			pc++
			
		default:
			log.Fatalf("Unknown op code at pc %v: %v (%v)", pc, op.Code(), op)
		}
	}

	return nil
}
