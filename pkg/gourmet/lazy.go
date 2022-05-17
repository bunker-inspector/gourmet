package gourmet

import (
	"sync"
)

func stream[T any](processor func(out chan T)) chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		processor(out)
	}()
	return out
}

func Seq[T any](in ...T) chan T {
	f := func(out chan T) {
		for _, v := range in {
			out <- v
		}
	}
	return stream(f)
}

func Map[T any, U any](fn func(T) U, in chan T) chan U {
	f := func(out chan U) {
		for v := range in {
			out <- fn(v)
		}
	}
	return stream(f)
}

func Each[T any](fn func(T), in chan T) chan T {
	f := func(out chan T) {
		for v := range in {
			fn(v)
		}
	}
	return stream(f)
}

func Tap[T any](fn func(T), in chan T) chan T {
	f := func(out chan T) {
		for v := range in {
			fn(v)
			out <- v
		}
	}
	return stream(f)
}

func Reduce[T any](reducer func(T, T) T,
	agg T,
	in chan T) chan T {
	f := func(out chan T) {
		for v := range in {
			agg = reducer(v, agg)
			out <- agg
		}
	}
	return stream(f)
}

func Filter[T any](pred func(T) bool, in chan T) chan T {
	f := func(out chan T) {
		for v := range in {
			if pred(v) {
				out <- v
			}
		}
	}
	return stream(f)
}

func Interleave[T any](ins ...chan T) chan T {
	f := func(out chan T) {
		buffer := make([]T, len(ins))
		for {
			for idx, in := range ins {
				received, ok := <-in
				if ok {
					buffer[idx] = received
				} else {
					return
				}
			}
			for _, sending := range buffer {
				out <- sending
			}

		}
	}
	return stream(f)
}

func Cycle[T any](in ...T) chan T {
	f := func(out chan T) {
		for true {
			for _, v := range in {
				out <- v
			}
		}
	}
	return stream(f)
}

func Iterate[T any](v T, fn func(T) T) chan T {
	f := func(out chan T) {
		for true {
			out <- v
			v = fn(v)
		}
	}
	return stream(f)
}

func Take[T any](n int, in chan T) chan T {
	f := func(out chan T) {
		i := 0
		for i < n {
			v, ok := <-in
			if !ok {
				close(in)
				break
			}
			out <- v
			i++
		}
	}
	return stream(f)
}

func TakeEvery[T any](n int, in chan T) chan T {
	f := func(out chan T) {
		v, ok := <-in
		for ok {
			out <- v
			for j := 0; j < n; j++ {
				v, ok = <-in
			}
		}
	}
	return stream(f)
}

func TakeWhile[T any](in chan T, pred func(T) bool) chan T {
	f := func(out chan T) {
		v, ok := <-in
		for ok && pred(v) {
			out <- v
			v, ok = <-in
		}
	}
	return stream(f)
}

func Fork[T any](in <-chan T) (chan T, chan T) {
	out := Distribute(in, 2)
	return out[0], out[1]
}

func Distribute[T any](in <-chan T, ct int) []chan T {
	var wg sync.WaitGroup

	out := make([]chan T, ct)
	for i := 0; i < ct; i++ {
		out[i] = make(chan T, 1000)
	}

	go func() {
		for v := range in {
			wg.Add(ct)
			for _, o := range out {
				go func(o chan T) {
					defer wg.Done()
					o <- v
				}(o)
			}
			wg.Wait()
		}
		for _, o := range out {
			close(o)
		}
	}()

	return out
}
