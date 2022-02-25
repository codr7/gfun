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
		case ADD:
			fmt.Printf("ADD %v %v\n", op.Reg1(), op.Reg2())
			v1 := &env.Regs[op.Reg1()]
			v2 := env.Regs[op.Reg2()]
			v1.Type().(NumType).AddVal(v1, v2)
			pc++
			break

		case DEC:
			fmt.Printf("DEC %v\n", op.Reg1())
			v := &env.Regs[op.Reg1()]
			d, err := v.Data()

			if err != nil {
				return err
			}
			
			v.Init(v.Type(), d.(int)-1)
			pc++
			break

		case INC:
			fmt.Printf("INC %v\n", op.Reg1())
			v := &env.Regs[op.Reg1()]
			d, err := v.Data()

			if err != nil {
				return err
			}
			
			v.Init(v.Type(), d.(int)+1)
			pc++
			break

		case SUB:
			fmt.Printf("SUB %v %v\n", op.Reg1(), op.Reg2())
			v1 := &env.Regs[op.Reg1()]
			v2 := env.Regs[op.Reg2()]
			v1.Type().(NumType).SubVal(v1, v2)
			pc++
			break

		case STOP:
			fmt.Printf("STOP\n")
			return nil

		default:
			log.Fatalf("Unknown op code at pc %v: %v (%v)", pc, op.Code(), op)
		}
	}

	return nil
}
