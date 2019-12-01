package lazy

import (
	"container/list"
	"sync"
)

func stream(processor func(out chan interface{})) chan interface{} {
	out := make(chan interface{})
	go processor(out)
	return out
}

func consume(processor func(out chan interface{}), in chan interface{}) chan interface{} {
	p := func(out chan interface{}) {
		processor(out)
		out <- nil
		close(in)
	}
	return stream(p)
}

func Seq(in ...interface{}) chan interface{} {
	s := func(out chan interface{}) {
		for _, v := range(in) {
			out <- v
		}
		out <- nil
	}
	return stream(s)
}

func Map(fn func(interface{})interface{}, in chan interface{}) chan interface{} {
	m := func(out chan interface{}) {
		v := <-in
		for v != nil {
			out <- fn(v)
			v = <-in
		}
	}
	return consume(m, in)
}

func Each(fn func(interface{}), in chan interface{}) chan interface{} {
	e := func(_ chan interface{}) {
		v := <-in
		for v != nil {
			fn(v)
			v = <-in
		}
	}
	return consume(e, in)
}

func Tap(fn func(interface{}), in chan interface{}) chan interface{} {
	t := func(out chan interface{}) {
		v := <-in
		for v != nil {
			fn(v)
			out <- v
			v = <-in
		}
	}
	return consume(t, in)
}


func Reduce(reducer func(interface{}, interface{})interface{},
	agg interface{},
	in chan interface{}) chan interface{} {
	r := func(out chan interface{}) {
		v := <-in
		for v != nil {
			agg = reducer(v, agg)
			out <- agg
			v = <-in
		}
	}
	return consume(r, in)
}

func Filter(pred func(interface{})bool, in chan interface{}) chan interface{} {
	f := func(out chan interface{}) {
		v := <-in
		for v != nil {
			if pred(v) {
				out <- v
			}
			v = <-in
		}
	}
	return consume(f, in)
}

func Interleave(ins ...chan interface{}) chan interface{} {
	z := func(out chan interface{}) {
		complete := 0
		for len(ins) > complete {
			for i, currIn := range(ins) {
				if currIn == nil {
					continue
				}
				nxt := <- currIn
				if nxt != nil {
					out <- nxt
				} else {
					close(currIn)
					ins[i] = nil
					complete++
				}
			}
		}
		out <- nil
	}
	return stream(z)
}

func Cycle(in ...interface{}) chan interface{} {
	c := func(out chan interface{}){
		for true {
			for _, v := range(in) {
				out <- v
			}
		}
	}
	return stream(c)
}

func Iterate(v interface{}, fn func(interface{}) interface{}) chan interface{} {
	i := func(out chan interface{}) {
		for true {
			out <- v
			v = fn(v)
		}
	}
	return stream(i)
}

func Take(n int, in chan interface{}) chan interface{} {
	t := func(out chan interface{}) {
		i := 0
		for i < n {
			v := <-in
			out <- v
			i++
			if v == nil {
				close(in)
				break
			}
		}
		out <- nil
	}
	return stream(t)
}

func TakeEvery(n int, in chan interface{}) chan interface{} {
	t := func(out chan interface{}) {
		v := <-in
		for v != nil {
			out <- v
			for j := 0; j < n && v != nil; j++ {
				v = <-in
			}
		}
		if v == nil {
			close(in)
		}
		out <- nil
	}
	return stream(t)
}

func TakeWhile(in chan interface{}, pred func(interface{}) bool) chan interface{} {
	t := func(out chan interface{}) {
		v := <-in
		for v != nil && pred(v) {
			out <- v
			v = <-in
		}
		if v == nil {
			close(in)
		}
		out <- nil
	}
	return stream(t)
}

func Fork(in <-chan interface{}) (a chan interface{}, b chan interface{}) {
	a = make(chan interface{})
	b = make(chan interface{})

	go func() {
		//buffer
		abuf := list.New()
		bbuf := list.New()

		//notifiers
		abufn := make(chan bool)
		bbufn := make(chan bool)

		notify := func(c chan bool) {c<-true}

		amut := &sync.Mutex{}
		bmut := &sync.Mutex{}

		processBuffer := func(buffer *list.List, output chan interface{}, n chan bool, m *sync.Mutex) {
			for true {
				//sleep until there's work to do
				<-n

				m.Lock()
				v := buffer.Front()
				buffer.Remove(v)
				m.Unlock()
				output <- v.Value

				//are we done?
				if v.Value == nil {
					return
				}
			}
		}

		go processBuffer(abuf, a, abufn, amut)
		go processBuffer(bbuf, b, bbufn, bmut)

		v := <-in
		for v != nil {
			amut.Lock()
			abuf.PushBack(v)
			amut.Unlock()
			go notify(abufn)

			bmut.Lock()
			bbuf.PushBack(v)
			bmut.Unlock()
			go notify(bbufn)

			v = <-in
		}
		amut.Lock()
		abuf.PushBack(nil)
		amut.Unlock()
		go notify(abufn)

		bmut.Lock()
		bbuf.PushBack(nil)
		bmut.Unlock()
		go notify(bbufn)
	}()
	return
}
