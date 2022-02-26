package gfun

import (
	"bufio"
	"fmt"
)

type ERead struct {
	E
	pos Pos
}

func NewERead(pos Pos, message string, args...interface{}) ERead {
	var self ERead
	self.Init(message, args...)
	self.pos = pos
	return self
}

func (self ERead) Error() string {
	return fmt.Sprintf("Error in %v: %v", self.pos, self.message)
}

var (
	defaultReaders = []Reader{ReadId, ReadLit}
)

func DefaultReaders() []Reader {
	return defaultReaders
}

type Reader func([]Reader, *bufio.Reader, *Pos, *M) (Form, error)

func ReadForm(readers []Reader, in *bufio.Reader, pos *Pos, m *M) (Form, error) {
	for _, r := range readers {
		if f, err := r(readers, in, pos, m); f != nil || err != nil {
			return f, err
		}
	}
	
	return nil, nil
}


func ReadId(readers []Reader, in *bufio.Reader, pos *Pos, m *M) (Form, error) {
	return nil, nil
}

func ReadLit(readers []Reader, in *bufio.Reader, pos *Pos, m *M) (Form, error) {
	return nil, nil
}
