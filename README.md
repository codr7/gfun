## gfun

```
fun: fib [n Int] [Int]
  (if (< n 2) n (fib (- n 1)))

(fib 10)

55
```

### intro
gfun aims to become a Lispy scripting language/vm implemented and embedded in Go.<br/><br/>

Much along the same lines as [g-fu](https://github.com/codr7/g-fu), but a tighter implementation that takes everything learned since into account. Among other changes; the VM now uses actual bytecode rather than structs and interfaces, and registers in place of stacks; and macro expansion works more like you would expect it to. The syntax is in many ways more traditional, but also simply different in places because of preferences picked up over time.<br/><br/>

The different layers of the implementation are well separated and may be used separately or mixed and matched; for example by putting a different syntax on top or adapting the existing one, or adding additional compilation steps..<br/><br/>

I intend to keep the implementation simple and small enough to be fun to play around with for educational purposes; which also ensures it stays reasonably general purpose and stable/bug free.

### limitations
For performance reasons, the evauluation loop specifies all operations inline, which means that it's impossible to extend the VM with new operations from the outside. Limits on number of operations, number of registers etc. are determined by the bytecode [format](https://github.com/codr7/gfun/blob/main/lib/op.go).

### status
It's still early days, I'm currently working my way towards calulating the Fibonacci sequence as above to get a comparable performance number on the implementation so far.

### types
#### Fun
#### Bool
The boolean type has two values, `T` and `F`.
#### Int
#### Nil
The nil type has one value, '_'.

### support
Should you wish to support this effort and allow me to spend more of my time and energy on evolving gfun, feel free to [help](https://liberapay.com/andreas7/donate) make that economically feasible.

### Go coder/mentor for hire
I'm currently available for hire.<br/>
Remote or relocation within Europe.<br/>
Send a message to codr7 at protonmail and I'll get back to you asap.
