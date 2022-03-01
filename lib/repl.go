package gfun

import (
	"bufio"
	"fmt"
	"io"
	"os"
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

			for _, f := range forms {
				if err := f.Emit(self); err != nil {
					fmt.Fprintln(out, err)
					break
				}
			}
			
			if self.emitPc != pc {
				self.DumpOps(pc, os.Stdout)
				self.EmitStop()
				self.Env().Regs[0].Init(&self.NilType, nil)
				
				if err := self.Eval(pc); err != nil {
					fmt.Fprintln(out, err)
				}
				
				resVal := self.Env().Regs[0]

				if resVal.Type() != &self.NilType {
					resVal.Type().DumpVal(resVal, out)
					fmt.Fprintf(out, "\n")
				}
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
