package eager

func TakeAll(seq <-chan interface{}) []interface{} {
	var result []interface{}
	v := <-seq
	for v != nil {
		result = append(result, v)
		v = <-seq
	}
	return result
}
