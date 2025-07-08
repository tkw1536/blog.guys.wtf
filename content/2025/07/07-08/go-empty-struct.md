---
title:          Using Empty Go Structs for a Hashmap
date:           2025-07-08
author:         Tom Wiesing 
authorLink:     https://tkw01536.de

description:    You should use empty go structs as the value type for a hashmap.

draft:          true
---

This morning I read a [post on the go blog](https://go.dev/blog/generic-interfaces) which eventually implemented a HashSet as:

```go
type HashSet[E comparable] map[E]bool

func (s HashSet[E]) Insert(v E)       { s[v] = true }
func (s HashSet[E]) Delete(v E)       { delete(s, v) }
func (s HashSet[E]) Has(v E) bool     { return s[v] }
func (s HashSet[E]) All() iter.Seq[E] { return maps.Keys(s) }
```

Having written a bunch of go myself, this immediately made me ask the question:

> Why didn't they use a `map[E]struct{}` ?

Knuth stated in his paper [^1] `We should forget about small efficiencies, say about 97% of the time: premature optimization is the root of all evil`.
So it this a "premature optimization", or is it one of the 3% where the optimization is worthwhile [^2]?

I asked a friend who works with a lot of go, and he pretty much said just that:

> Unless your set needs to host millions of items do you really need to save bools worth of space?
> 
> To me it always felt like an unnecessary optimization 

Let us investigate this by actually actually implementing a struct-based version:

```go
type HashSetStruct[E comparable] map[E]struct{}

func (s HashSetStruct[E]) Insert(v E)       { s[v] = struct{}{} }
func (s HashSetStruct[E]) Delete(v E)       { delete(s, v) }
func (s HashSetStruct[E]) Has(v E) bool     { _, ok := s[v]; return ok }
func (s HashSetStruct[E]) All() iter.Seq[E] { return maps.Keys(s) }
```

And making a benchmark like:

```go
// some pseudo-random integers for testing
var ints [10_000]int

func Benchmark_HashSetStruct(b *testing.B) {
	for b.Loop() {
		set := make(HashSetStruct[int], len(ints))
		for _, v := range ints {
			set.Insert(v)
		}
		for _, v := range ints {
			_ = set.Has(v)
		}
	}
}
```

Writing a second benchmark for the plain `HashSet` and running both benchmarks [^3] gave [^4]:

```
$ go test -bench . -benchmem
Benchmark_HashSet-10                8763            136473 ns/op          295554 B/op         33 allocs/op
Benchmark_HashSetStruct-10          8548            134898 ns/op          295554 B/op         33 allocs/op
```

This doesn't seem to show any significant performance difference - even the same number of allocations in both cases. 
There are only a few nanoseconds of difference - none in number of allocations or bytes used. 
So does this mean that it's not worth it?

For one, this benchmark is rather synthetic and only tests insertion and lookup.
Even leaving alone this particular benchmark it can be said that since go is a compiled language, optimizations should happen under the hood. 
If the compiler could optimize both cases, it would make it mostly a stylistic choice if you want to use `struct{}` or `bool`. 
Specifically in this case both implementation expose the same methods - so here it doesn't matter. 


But how could the compiler optimize the `bool` case?
It would have to prove that the `false` value is never used. 
Consider this snippet of code:

```go
set := make(HashSet[int])
set.Insert(42)
// ... sometime later ...
set[69] = false
```

Here the value of `false` is explicitly used. 
Because the HashSet struct does not hide its' implementation in an unexported field this is perfectly legitimate[^5]. 
What happens in such a case?

For one, the implementation of the `All` method would be wrong, as it would still return both keys. 
It does not check if keys in the map are `true`, but blindly returns all keys.
A correct implementation of the `All()` method requires checking if the associated value is true:

```go
func (s HashSet[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		for k, v := range s {
			if !v {
				continue
			}
			if !yield(k) {
				return
			}
		}
	}
}
```

This in turn means it is unsound for the compiler to optimize the value of `false` away, as there is an observable difference if false is written to the map.

Summing it all up:

- My benchmarks didn't show any significant performance difference between using `map[E]struct{}` and `map[E]bool`.
- It is a stylistic choice what you want to use, but only `map[E]struct{}` can allow the compiler to optimize.
- The implementation of the `All()` method is wrong.

[^1]: Donald E. Knuth. 1974. Structured Programming with go to Statements. ACM Comput. Surv. 6, 4 (Dec. 1974), 261–301. doi: [10.1145/356635.356640](https://doi.org/10.1145/356635.356640).

[^2]: An interesting sidenote from Knuth's paper is that it is frequently misinterpreted, not all optimization is pre-mature. 
See for example the blog post ["Revisiting Knuth’s “Premature Optimization” Paper"](https://web.archive.org/web/20250619231836/https://probablydance.com/2025/06/19/revisiting-knuths-premature-optimization-paper/) by *Malte Skarupke*.

[^3]: Full code available [GitHub](https://gist.github.com/tkw1536/f3a6f89f9c49a36f6143a426014630cb). 

[^4]: Part of the benchmark output omitted for brevity. 

[^5]: Another means of making this illegitimate might be to write human-readable documentation that explicitly forbids writing to the underlying map directly.
But this could not be understood by the compiler, and also isn't being done here. 

[^6]: A correct implementation of the `All` method might be 