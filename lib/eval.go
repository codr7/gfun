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
		case BRANCH_REG:
			fmt.Printf("BRANCH_REG %v %v\n", op.Reg1())
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

		case DEC:
			fmt.Printf("DEC %v %v\n", op.Reg1(), op.Reg2())
			l := &env.Regs[op.Reg1()]
			r := env.Regs[op.Reg2()]
			var lv interface{}
			var err error
			
			if lv, err = l.Data(); err != nil {
				return err
			}
			
			var rv interface{}
			
			if rv, err = r.Data(); err != nil {
				return err
			}
			
			l.Init(l.Type(), lv.(int)-rv.(int))
			pc++

		case INC:
			fmt.Printf("INC %v %v\n", op.Reg1(), op.Reg2())
			l := &env.Regs[op.Reg1()]
			r := env.Regs[op.Reg2()]
			var lv interface{}
			var err error
			
			if lv, err = l.Data(); err != nil {
				return err
			}
			
			var rv interface{}
			
			if rv, err = r.Data(); err != nil {
				return err
			}
			
			l.Init(l.Type(), lv.(int)+rv.(int))
			pc++

		case STOP:
			fmt.Printf("STOP\n")
			return nil

		default:
			log.Fatalf("Unknown op code at pc %v: %v (%v)", pc, op.Code(), op)
		}
	}

	return nil
}
