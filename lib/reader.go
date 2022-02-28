package gfun

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
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
	defaultReaders = []Reader{ReadWs, ReadCall, ReadInt, ReadId}
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

func ReadCall(readers []Reader, in *bufio.Reader, pos *Pos, m *M) (Form, error) {
	fpos := *pos
	var c rune
	
	if c, _, _ = in.ReadRune(); c == '(' {
		pos.Col++
	} else {
		in.UnreadRune()
		return nil, nil
	}

	var t Form
	var err error

	if t, err = ReadForm(readers, in, pos, m); err != nil {
		return nil, err
	}
	
	var as []Form

	for {
		if a, err := ReadForm(readers, in, pos, m); err != nil {
			return nil, err
		} else if a == nil {
			break
		} else {
			as = append(as, a)
		}
	}

	if c, _, _ = in.ReadRune(); c != ')' {
		return nil, NewERead(fpos, "Open call")
	}

	pos.Col++
	return NewCallForm(fpos, t, as), nil
}

func ReadId(readers []Reader, in *bufio.Reader, pos *Pos, m *M) (Form, error) {
	fpos := *pos
	var buf strings.Builder

	for {
		if c, _, err := in.ReadRune(); err == io.EOF {
			in.UnreadRune()
			break
		} else if err != nil {
			return nil, err
		} else if unicode.IsSpace(c) ||
			c == '(' || c == ')' || c == '[' || c == ']' ||
			(buf.Len() > 0 && c == '|') {
			in.UnreadRune()
			break
		} else {
			pos.Col++
			buf.WriteRune(c)
		}
	}

	if buf.Len() == 0 {
		return nil, nil
	}
	
	return NewIdForm(fpos, m.Sym(buf.String())), nil
}

func ReadInt(readers []Reader, in *bufio.Reader, pos *Pos, m *M) (Form, error) {
	fpos := *pos
	var buf strings.Builder

	for {
		c, _, e := in.ReadRune()

		if e == io.EOF {
			break
		}
		
		if e != nil {
			return nil, e
		}

		if c != '-' && !unicode.IsDigit(c) {
			in.UnreadRune()
			break
		}

		buf.WriteRune(c)
		pos.Col++
	}

	if buf.Len() == 0 {
		return nil, nil
	}

	s := buf.String()

	if s == "-" {
		return NewIdForm(fpos, m.Sym(s)), nil
	}
	
	n, e := strconv.ParseInt(s, 10, 64)

	if e != nil {
		return nil, NewERead(fpos, "Invalid Int: %v", s)
	}

	return NewLitForm(fpos, NewVal(&m.IntType, int(n))), nil
}

func ReadWs(readers []Reader, in *bufio.Reader, pos *Pos, m *M) (Form, error) {
	for {
		if c, _, err := in.ReadRune(); err != nil {
			return nil, err
		} else if unicode.IsSpace(c) {
			if c == '\n' {
				pos.Row++
				pos.Col = 0
			} else {
				pos.Col++
			}
		} else {
			in.UnreadRune()
			break
		}
	}
	
	return nil, nil
}
