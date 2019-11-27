package eager

func Collect(seq chan interface{}) []interface{} {
	var result []interface{}
	v := <-seq
	for v != nil {
		result = append(result, v)
		v = <-seq
	}
	close(seq)
	return result
}
