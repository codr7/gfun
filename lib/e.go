package gfun

import (
	"fmt"
)

type E struct {
	message string
}

func (self *E) Init(message string, args...interface{}) *E {
	self.message = fmt.Sprintf(message, args...)
	return self
}

func (self E) Error() string {
	return fmt.Sprintf("Error: %v", self.message)
}
