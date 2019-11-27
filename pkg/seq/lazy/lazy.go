package lazy

func Seq(o ...interface{}) chan interface{} {
	c := make(chan interface{})
	go func(){
		for _, v := range(o) {
			c <- v
		}
		c <- nil
		close(c)
	}()
	return c
}

func Map(fn func(interface{})interface{}, seq <-chan interface{}) chan interface{} {
	c := make(chan interface{})
	go func(){
		v := <-seq
		for v != nil {
			c <- fn(v)
			v = <-seq
		}
		c <- nil
	}()
	return c
}

func Filter(pred func(interface{})bool, seq <-chan interface{}) chan interface{} {
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
	}()
	return c
}

func Reverse(o <-chan interface{}) {
	panic("Oh god oh god oh shit oh god")
}

func Zip(os ...<-chan interface{}) chan interface{} {
	c := make(chan interface{})
	go func() {
		for len(os) > 0 {
			for i, o := range(os) {
				nxt := <- o
				if nxt != nil {
					c <- nxt
				} else {
					if i+1 >= len(os) {
						os = os[:i]
					} else {
						os = append(os[:i], os[i+1:]...)
					}
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
		for i < n {
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
		for pred(v) {
			c <- v
			v = <-seq
		}
		c <- nil
	}()
	return c
}
