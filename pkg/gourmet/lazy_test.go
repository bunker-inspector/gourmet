package gourmet

import (
	"testing"
	"strconv"
)

func inc (i int) int {
	return i + 1
}

func dec (i int) int {
	return i - 1
}


func TestCycle(t *testing.T) {
	s := Cycle(1,2,3)
	d := Take(4, s)
	for _, v := range([]int{1,2,3,1}) {
		curr := <-d
		if v != curr {
			t.Fatalf("Cycle failed: Expected %d, got %d", v, curr)
		}
	}
}

func lessThan(n int) func(int)bool {
	return func(i int) bool {
		return i < n
	}
}


func TestTake(t *testing.T) {
	s := Cycle(1,2,3)
	f := Collect(Take(3, s))
	for i, v := range([]int{1,2,3}) {
		if v != f[i] {
			t.Fatalf("Take failed: Expected %d, got %d", v, f[i])
		}
	}
}

func TestTakeWhile(t *testing.T) {
	s := Seq(1,2,3,4,5)
	e := Collect(TakeWhile(s, lessThan(4)))
	for i, v := range([]int{1,2,3}) {
		if v != e[i] {
			t.Fatalf("TakeWhile failed: Expected %d, got %d", v, e[i])
		}
	}
}

func TestTakeEvery(t *testing.T) {
	s := Seq(1,2,3,4,5,6)
	e := Collect(TakeEvery(2, s))
	for i, v := range([]int{1,3,5}) {
		if v != e[i] {
			t.Fatalf("TakeEvery failed: Expected %d, got %d", v, e[i])
		}
	}
}

func TestIterate(t *testing.T) {
	inc := func(i int) int {return i + 1}
	s := Take(4, Iterate(1, inc))
	for _, v := range([]int{1,2,3,4}) {
		curr := <-s
		if v != curr {
			t.Fatalf("Iterate failed: Expected %d, got %d", v, curr)
		}
	}
}

func TestMap(t *testing.T) {
	s := Collect(Map(inc, Seq(1,2,3,4)))
	for i, v := range([]int{2,3,4,5}) {
		if v != s[i] {
			t.Fatalf("Map failed: Expected %d, got %d", v, s[i])
		}
	}
}

func TestMapConverts(t *testing.T) {
	s := Collect(Map(strconv.Itoa, Seq(1,2,3,4)))
	for i, v := range([]string{"1","2","3","4"}) {
		if v != s[i] {
			t.Fatalf("Map failed: Expected %s, got %s", v, s[i])
		}
	}
}

func TestReduce(t *testing.T) {
	reducer := func(v int, sum int) int {
		return v + sum
	}
	s := Collect(Reduce(reducer, 0, Seq(1,2,3)))
	sum := s[len(s)-1]
	if sum != 6 {
		t.Fatalf("Reduce failed: Expected %d, got %d", 6, s[len(s)-1])
	}
}

func TestEach(t *testing.T) {
	s := Seq(1,2,3)
	r := []int{}
	e := func(i int) {
		r = append(r, i)
	}
	<-Each(e, s)
	for i, v := range([]int{1,2,3}) {
		if v != r[i] {
			t.Fatalf("Each failed: Expected %d, got %d", v, r[i])
		}
	}
}

func isEven(i int) bool {
	return i % 2 == 0
}

func TestFilter(t *testing.T) {
	a := Seq(1,2,3,4)
	b := Collect(Filter(isEven, a))
	for i, v := range([]int{2,4}) {
		if v != b[i] {
			t.Fatalf("Filter failed: Expected %d, got %d", v, b[i])
		}
	}
}

func TestInterleave(t *testing.T) {
	a := Seq(1,3,5)
	b := Seq(2,4,6)
	c := Collect(Interleave(a,b))
	for i, v := range([]int{1,2,3,4,5,6}) {
		if v != c[i] {
			t.Logf("Interleave failed: Expected %d, got %d at step %d", v, c[i], i)
			t.Fail()
		}
	}
}

func TestDistribute(t *testing.T) {
	s := Seq(2,3,4,5)
	out := Distribute(s, 3)
	ar := Collect(Map(inc, out[0]))
	br := Collect(Map(dec, out[1]))
	cr := Collect(Map(dec, out[2]))

	for i, v := range([]int{2,3,4,5}) {
		if v+1 != ar[i] {
			t.Fatalf("Distribute failed: Expected %d from ar, got %d", v+1, ar[i])
		}
		if v-1 != br[i] {
			t.Fatalf("Distribute failed: Expected %d from br, got %d", v-1, br[i])
		}
		if v-1 != cr[i] {
			t.Fatalf("Distribute failed: Expected %d from br, got %d", v-1, cr[i])
		}
	}
}
