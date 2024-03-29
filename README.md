## Gourmet
A collection of high level functions for stream processing using channels to model laziness

#### A few examples

##### Reduce
```go
reducer := func(v int, sum int) int {
    return v + sum
}
s := Collect(Reduce(reducer, 0, Seq(1,2,3)))
sum := s[len(s)-1]
 //sum: 6
```

##### Each
```go
s := Seq(1,2,3)
r := []int{}
e := func(i int) {
	r = append(r, i)
}
<-Each(e, s)
 //r: [1,2,3]
```

##### Interleave
```go
a := Seq(1,3,5)
b := Seq(2,4,6)
c := Collect(Interleave(a,b))
//c: [1,2,3,4,5,6]
```
