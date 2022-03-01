## gfun

```
(fun: fib [n Int] Int
  (if (< n 2) n (+ (fib (- n 1)) (fib (- n 2)))))

(fib 10)

55
```

### intro
gfun aims to become a Lispy scripting language/vm implemented and embedded in Go.<br/><br/>

Much along the same lines as [g-fu](https://github.com/codr7/g-fu), but a tighter implementation that takes everything learned since into account.<br/><br/>

The different layers of the implementation are well separated and may be used separately or mixed and matched; for example by putting a different syntax on top or adapting the existing one, or adding additional compilation steps.<br/><br/>

### limitations
For performance reasons, the core loop specifies operations inline; which means that it's impossible to extend with new operations without changing the code. Limits on number of operations, number of registers etc. are determined by the bytecode [format](https://github.com/codr7/gfun/blob/main/lib/op.go).

### status
It's still early days, I'm currently profiling and optimizing the implementation guided by initial performance numbers.

### types
#### Any
The root type.
#### Fun < Any
The type of functions.
#### Bool < Any
The boolean type has two values, `T` and `F`.
#### Int < Any
The type of whole numbers.
#### Macro < Any
The type of macros.
#### Meta < Any
The type of types.
#### Nil < Any
The nil type has one value, `_`.

### performance
gfun is pretty slow at the moment, around 100 times slower than Python; I know it is possible to go a lot faster, the remaining task is figuring out what part of the design is causing trouble.

```
(fun: fib [n Int] Int
  (if (< n 2) n (+ (fib (- n 1)) (fib (- n 2)))))

(bench 100 (fib 20))

9580
```

### support
Should you wish to support this effort and allow me to spend more of my time and energy on evolving gfun, feel free to [help](https://liberapay.com/andreas7/donate) make that economically feasible.

### coder/mentor for hire
I'm currently available for hire.<br/>
Remote or relocation within Europe.<br/>
Send a message to codr7 at protonmail and I'll get back to you asap.
