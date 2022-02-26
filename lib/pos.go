package gfun

type Pos struct {
	source string
	Row, Col int
}

func NewPos(source string, row, col int) Pos {
	var self Pos
	self.Init(source, row, col)
	return self
}

func (self *Pos) Init(source string, row, col int) {
	self.source = source
	self.Row = row
	self.Col = col
}
