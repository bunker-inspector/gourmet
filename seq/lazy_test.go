package seq

import (
	//"fmt"
	"testing"
)

func TestCycle(t *testing.T) {
	s := Cycle([]interface{}{1,2,3}...)
	d := Take(4, s)
	for i, v := range([]int{1,2,3,1}) {
		if v != d[i] {
			t.Fatalf("Cycle failed: Expected %d, go %d", v, d[i])
		}
	}
}

func TestIterate(t *testing.T) {
	inc := func(i interface{}) interface{} {return i.(int) + 1}
	s := Take(4, Iterate(1, inc))
	for i, v := range([]int{1,2,3,4}) {
		if v != s[i] {
			t.Fatalf("Iterate failed: Expected %d, go %d", v, s[i])
		}
	}
}

func TestMap(t *testing.T) {
	inc := func(i interface{}) interface{} {return i.(int) + 1}
	s := Take(4, Map(inc, LazySeq(1,2,3,4)))
	for i, v := range([]int{2,3,4,5}) {
		if v != s[i] {
			t.Fatalf("Iterate failed: Expected %d, go %d", v, s[i])
		}
	}
}

func TestZip(t *testing.T) {
	a := LazySeq(1,3,5)
	b := LazySeq(2,4,6)
	c := Take(6, Zip(a,b))
	for i, v := range([]int{1,2,3,4,5,6}) {
		if v != c[i] {
			t.Fatalf("Zip failed: Expected %d, go %d", v, c[i])
		}
	}
}
