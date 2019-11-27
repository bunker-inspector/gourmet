package lazy

func Seq(o ...interface{}) chan interface{} {
	c := make(chan interface{})
	go func(){
		for _, v := range(o) {
			c <- v
		}
		c <- nil
	}()
	return c
}

func Map(fn func(interface{})interface{}, seq chan interface{}) chan interface{} {
	c := make(chan interface{})
	go func(){
		v := <-seq
		for v != nil {
			c <- fn(v)
			v = <-seq
		}
		c <- nil
		close(seq)
	}()
	return c
}

func Filter(pred func(interface{})bool, seq chan interface{}) chan interface{} {
	c := make(chan interface{})
	go func(){
		v := <-seq
		for v != nil {
			if pred(v) {
				c <- v
			}
			v = <-seq
		}
		c <- nil
		close(seq)
	}()
	return c
}

func Reverse(o <-chan interface{}) {
	panic("Oh god oh god oh shit oh god")
}

func Zip(os ...chan interface{}) chan interface{} {
	c := make(chan interface{})
	go func() {
		complete := 0
		for len(os) > complete {
			for i, o := range(os) {
				if o == nil {
					continue
				}
				nxt := <- o
				if nxt != nil {
					c <- nxt
				} else {
					close(o)
					os[i] = nil
					complete++
				}
			}
		}
		c <- nil
	}()
	return c
}

func Cycle(o ...interface{}) chan interface{} {
	c := make(chan interface{})
	go func(){
		for true {
			for _, v := range(o) {
				c <- v
			}
		}
	}()
	return c
}

func Iterate(v interface{}, fn func(interface{}) interface{}) chan interface{} {
	c := make(chan interface{})
	go func() {
		for true {
			c <- v
			v = fn(v)
		}
	}()
	return c
}

func Take(n int, seq <-chan interface{}) chan interface{} {
	c := make(chan interface{})
	v := <-seq
	go func() {
		i := 0
		for i < n && v != nil {
			c <- v
			v = <-seq
			i++
		}
		c <- nil
	}()
	return c
}

func TakeWhile(seq <-chan interface{}, pred func(interface{}) bool) chan interface{} {
	c := make(chan interface{})
	v := <-seq
	go func() {
		for v != nil && pred(v) {
			c <- v
			v = <-seq
		}
		c <- nil
	}()
	return c
}

func Fork(seq <-chan interface{}) (a chan interface{}, b chan interface{}) {
	a = make(chan interface{})
	b = make(chan interface{})
	go func() {
		v := <-seq
		for v != nil {
			a <- v
			b <- v
			v = <-seq
		}
		a <- nil
		b <- nil
	}()
	return
}
