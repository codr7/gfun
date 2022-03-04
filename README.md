## GFun

```
  (fun: fib [n Int] Int
    (if (< n 2) n (+ (fib (dec n)) (fib (dec n 2)))))
  
  (fib 10)

55
```

### intro
GFun aims to become a Lispy scripting language/vm implemented and embedded in Go.<br/><br/>

Much along the same lines as [g-fu](https://github.com/codr7/g-fu), but a tighter implementation that takes everything learned since into account.<br/><br/>

The different layers of the implementation are well separated and may be used separately or mixed and matched; for example by putting a different syntax on top or adapting the existing one, or adding additional compilation steps.<br/><br/>

### limitations
For performance reasons, the core loop specifies operations inline; which means that it's impossible to extend with new operations without changing the code. Limits on number of operations, number of registers etc. are determined by the bytecode [format](https://github.com/codr7/gfun/blob/main/lib/op.go).

### status
It's still early days, I'm currently profiling and optimizing the implementation guided by initial performance numbers. All functionality described in this document is expected to work.

### setup

```
$ cd test
$ go test
$ cd ..
$ go build gfun.go
$ ./gfun test/all.gf
$ ./gfun
```

### types
GFun supports first class types, `typeof` may be used to get the type of any value.

```
  (typeof 42)

Int
  (typeof Int)

Meta
  (typeof Meta)

Meta
```

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

### bindings
Values may be temprarily bound to identifiers using `let`.

```
  (let [foo 35 bar (+ foo 7)] bar)

42
```

Values may be bound until further using `set`.

```
  (set foo 35 bar 7 baz (+ 35 7))
  baz

42
```

### functions
Functions may be anonymous.

```
  (let [foo (fun [bar Int] Int bar)]
    (dump foo)
    (foo 42))
  
(Fun 0x824634418592)
42
```

### performance
GFun is currently around 1-5 times as slow as Python.

```
  (fun: fibrec [n Int] Int
    (if (< n 2) n (+ (fibrec (dec n)) (fibrec (dec n 2)))))

  (bench 100 (fibrec 20))

562
```

```
  (fun: fibtail [n Int a Int b Int] Int
    (if (= n 0) a (if (= n 1) b (fibtail (dec n) b (+ a b)))))

  (bench 10000 (fibtail 70 0 1))

214
```

#### fusing
The generated bytecode is analyzed before evaluation in an attempt to fuse as much of it as possible.
In the following example; GFun initially detects that the arguments are not used, which results in the function entry point being moved to 3; then the exit point is moved past the `_` since that's the default return value.

```
 (fun: foo [x Int y Int] Nil _)

2022/03/02 22:13:46 Fused unused load at 1: 32
2022/03/02 22:13:46 Fused unused load at 2: 33
2022/03/02 22:13:46 Fused entry to 2
2022/03/02 22:13:46 Fused entry to 3
2022/03/02 22:13:46 Fused exit to 3
0 GOTO 5
1 NOP
2 NOP
3 RET
4 RET
```
### support
Should you wish to support this effort and allow me to spend more of my time and energy on evolving GFun, feel free to [help](https://liberapay.com/andreas7/donate) make that economically feasible.

### coder/mentor for hire
I'm currently available for hire.<br/>
Remote or relocation within Europe.<br/>
Send a message to codr7 at protonmail and I'll get back to you asap.
