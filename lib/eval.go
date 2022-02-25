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
			fmt.Printf("Stop\n")
			return nil
		case INC:
			fmt.Printf("Inc %v\n", op.Reg())
			val := &env.Regs[op.Reg()]
			d, err := val.Data()

			if err != nil {
				return err
			}
			
			val.Init(val.Type(), d.(int)+1)
			pc++
			break
		default:
			log.Fatalf("Unknown op code at pc %v: %v (%v)", pc, op.Code(), op)
		}
	}

	return nil
}
