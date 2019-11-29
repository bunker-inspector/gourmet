Gourmet is an Elixir standard library inspired collection of high level functions for stream processing using channels to model lazy sequences

#### A few examples

##### Reduce
```go
reducer := func(v interface{}, sum interface{}) interface{}{
    return v.(int) + sum.(int)
}
s := eager.Collect(Reduce(reducer, 0, Seq(1,2,3)))
sum := s[len(s)-1]
 //sum: 6
```

##### Each
```go
s := Seq(1,2,3)
r := []int{}
e := func(i interface{}) {
	r = append(r, i.(int))
}
<-Each(e, s)
 //r: [1, 2, 3]
```

##### Zip
```go
a := Seq(1,3,5)
b := Seq(2,4,6)
c := eager.Collect(Zip(a,b))
//c: [1,2,3,4,5,6]
```
