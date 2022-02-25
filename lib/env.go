package gfun

const (
	RegCount = 64
	ArgCount = 8
	RetCount = 4
)

type Env struct {
	outer *Env

	Regs [RegCount]Val
	Args [ArgCount]Val
	Rets [RetCount]Val
}

func (self *Env) Init(outer *Env) {
	self.outer = outer
}
