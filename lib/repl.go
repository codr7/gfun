package gfun

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func (self *M) Repl(readers []Reader, in io.Reader, out io.Writer) error {
	fmt.Fprintf(out, "  ")
	var buf strings.Builder
	ins := bufio.NewScanner(in)
	
	for ins.Scan() {
		if line := ins.Text(); len(line) == 0 && buf.Len() > 0 {
			bin := bufio.NewReader(strings.NewReader(buf.String()))
			pos := NewPos("repl", 0, 0)
			var forms []Form
			
			for {
				if f, err := ReadForm(readers, bin, &pos, self); err == io.EOF {
					break
				} else if err != nil {
					fmt.Fprintln(out, err)
					forms = nil
					break
				} else if f == nil {
					break
				} else {
					forms = append(forms, f)
				}
			}

			pc := self.emitPc
			
			for len(forms) != 0 {
				f := forms[0]
				var err error
				forms, err = f.Emit(forms[1:], self)
				
				if err != nil {
					fmt.Fprintln(out, err)
					break
				}
			}

			if len(forms) == 0 && self.emitPc != pc {
				self.EmitStop()
				self.env.Regs[0].Init(&self.NilType, nil)
				
				if err := self.Eval(pc); err != nil {
					fmt.Fprintln(out, err)
				}
				
				var err error
				var res interface{}
				
				if res, err = self.env.Regs[0].Data(); err != nil {
					return err
				}
				
				fmt.Fprintf(out, "%v\n", res)
			}
			
			buf.Reset()
		} else {
			buf.WriteString(line)
			buf.WriteRune('\n')
		}

		fmt.Fprintf(out, "  ")
	}

	return nil
}
