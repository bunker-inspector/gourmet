package lazy

import (
	"container/list"
	"sync"
)

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

func Take(n int, seq chan interface{}) chan interface{} {
	c := make(chan interface{})
	v := <-seq
	go func() {
		i := 0
		for i < n {
			c <- v
			v = <-seq
			i++
			if v == nil {
				close(seq)
				break
			}
		}
		c <- nil
	}()
	return c
}

func TakeEvery(n int, seq chan interface{}) chan interface{} {
	c := make(chan interface{})
	v := <-seq
	go func() {
		for v != nil {
			c <- v
			for j := 0; j < n && v != nil; j++ {
				v = <-seq
			}
			if v == nil {
				close(seq)
				break
			}
		}
		c <- nil
	}()
	return c
}

func TakeWhile(seq chan interface{}, pred func(interface{}) bool) chan interface{} {
	c := make(chan interface{})
	v := <-seq
	go func() {
		for v != nil && pred(v) {
			c <- v
			v = <-seq
		}
		c <- nil
		close(seq)
	}()
	return c
}

func Fork(seq <-chan interface{}) (a chan interface{}, b chan interface{}) {
	a = make(chan interface{})
	b = make(chan interface{})

	go func() {
		aBuffer := list.New()
		bBuffer := list.New()

		aBufferN := make(chan bool)
		bBufferN := make(chan bool)

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

		go processBuffer(aBuffer, a, aBufferN, amut)
		go processBuffer(bBuffer, b, bBufferN, bmut)

		v := <-seq
		for v != nil {
			amut.Lock()
			aBuffer.PushBack(v)
			amut.Unlock()
			go notify(aBufferN)

			bmut.Lock()
			bBuffer.PushBack(v)
			bmut.Unlock()
			go notify(bBufferN)

			v = <-seq
		}
		amut.Lock()
		aBuffer.PushBack(nil)
		amut.Unlock()
		go notify(aBufferN)

		bmut.Lock()
		bBuffer.PushBack(nil)
		bmut.Unlock()
		go notify(bBufferN)
	}()
	return
}
