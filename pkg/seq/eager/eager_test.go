package eager

import (
	"testing"
	"github.com/tedgkassen/gourmet/pkg/seq/lazy"
)

func TestTakeAll(t *testing.T) {
	s := lazy.LazySeq(1,2,3)
	d := TakeAll(s)
	for i, v := range([]int{1,2,3}) {
		if v != d[i] {
			t.Fatalf("TakeAll failed: Expected %d, go %d", v, d[i])
		}
	}
}
