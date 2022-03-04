(test Int
  (typeof 42))

(test Meta
  (typeof Int))

(test Meta
  (typeof Meta))
  
(test 42
  (let [foo 35 bar (+ foo 7)] bar))

(test 42
  (let [foo 42]
    (fun: bar [] Int foo)
    (bar)))

(test 3
  (let [foo 1]
    (fun: bar [] Int foo)
    (let [foo 3]
      (bar))))

(test 42
  (set foo 35 bar 7 baz (+ 35 7))
  baz)

(test 42
  (let [foo (fun [bar Int] Int bar)]
    (foo 42)))

(test 55
  (fun: fibrec [n Int] Int
    (if (< n 2) n (+ (fibrec (dec n)) (fibrec (dec n 2)))))

  (fibrec 10))

(test 55
  (fun: fibtail [n Int a Int b Int] Int
    (if (= n 0) a (if (= n 1) b (fibtail (dec n) b (+ a b)))))

  (fibtail 10 0 1))