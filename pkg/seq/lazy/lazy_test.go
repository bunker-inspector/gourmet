package lazy

import (
	//"fmt"
	"testing"
	"github.com/tedgkassen/gourmet/pkg/seq/eager"
)

func TestCycle(t *testing.T) {
	s := Cycle([]interface{}{1,2,3}...)
	d := Take(4, s)
	for _, v := range([]int{1,2,3,1}) {
		curr := <-d
		if v != curr {
			t.Fatalf("Cycle failed: Expected %d, got %d", v, curr)
		}
	}
}

func TestIterate(t *testing.T) {
	inc := func(i interface{}) interface{} {return i.(int) + 1}
	s := Take(4, Iterate(1, inc))
	for _, v := range([]int{1,2,3,4}) {
		curr := <-s
		if v != curr {
			t.Fatalf("Iterate failed: Expected %d, got %d", v, curr)
		}
	}
}

func TestMap(t *testing.T) {
	inc := func(i interface{}) interface{} {return i.(int) + 1}
	s := eager.Collect(Map(inc, Seq(1,2,3,4)))
	for i, v := range([]int{2,3,4,5}) {
		if v != s[i] {
			t.Fatalf("Iterate failed: Expected %d, got %d", v, s[i])
		}
	}
}

func TestZip(t *testing.T) {
	a := Seq(1,3,5)
	b := Seq(2,4,6)
	c := eager.Collect(Zip(a,b))
	for i, v := range([]int{1,2,3,4,5,6}) {
		if v != c[i] {
			t.Fatalf("Zip failed: Expected %d, got %d", v, c[i])
		}
	}
}
