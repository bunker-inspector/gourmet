package eager

import (
	"testing"
	"github.com/tedgkassen/gourmet/pkg/seq/lazy"
)

func TestCollect(t *testing.T) {
	s := lazy.Seq(1,2,3)
	d := Collect(s)
	for i, v := range([]int{1,2,3}) {
		if v != d[i] {
			t.Fatalf("TakeAll failed: Expected %d, go %d", v, d[i])
		}
	}
}
