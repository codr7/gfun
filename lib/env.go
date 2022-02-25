package gfun

const (
	RegCount = 64
	ArgCount = 8
	RetCount = 4
)

type Env struct {	
	Regs [RegCount]Val
	Args [ArgCount]Val
	Rets [RetCount]Val

	outer *Env
	bindings map[*Sym]int
}

func (self *Env) Init(outer *Env) {
	self.outer = outer
	self.bindings = make(map[*Sym]int)
}
