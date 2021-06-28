# Async

## Introduction

Async wraps golang library methods and provides a way to call them in a non blocking way. The context parameter adds more control over the lifecycle of the goroutine.

Features:
1. Context Compatible
2. All goroutines are immediately canceled if any one goroutine throws an error

## ErrorHandle

```go
type ErrorHandle func() error
```

The ErrorHandle function lets you handle the error coming from the async call, if any, in either of the 2 ways:
1. Call this method before you need the result for synchronous behaviour.
2. Defer it till the end of the caller lifecycle resulting in asynchronous behaviour.

## Basic Usage

```go
go get github.com/siddhant2408/async
```

### Copy

```go
// asynchronous
asyncErr := async.Copy(ctx, async.TeeReader{dst, src})
defer func(asyncErr async.ErrorHandle) {
	err := asyncErr()
	if err != nil {
		//handle error
	}
}(asyncErr)

// synchronous
asyncErr := async.CopyMultiple(ctx, []async.TeeReader{{dst1, src1}, {dst2, src2}})
err := asyncErr()
if err != nil {
	//handle error
}
```

## Results

The below results state that `io.Copy` shines for small copies whereas the `async.CopyMultiple` shines for large ones.

`go test -run=XXX -bench=. -benchmem -benchtime=10s`

```
BenchmarkCopyMultiple/sampleSize=32;numCopies=5-8         	 2735391	      4324 ns/op	    1800 B/op	      32 allocs/op
BenchmarkCopyMultiple/sampleSize=32;numCopies=10-8        	 1854004	      6458 ns/op	    3288 B/op	      53 allocs/op
BenchmarkCopyMultiple/sampleSize=32;numCopies=15-8        	 1402531	      8526 ns/op	    4648 B/op	      73 allocs/op
BenchmarkCopyMultiple/sampleSize=32;numCopies=20-8        	 1000000	     10722 ns/op	    6265 B/op	      94 allocs/op
BenchmarkCopyMultiple/sampleSize=1024;numCopies=5-8       	 1000000	     10562 ns/op	   21640 B/op	      32 allocs/op
BenchmarkCopyMultiple/sampleSize=1024;numCopies=10-8      	  682428	     17308 ns/op	   42969 B/op	      53 allocs/op
BenchmarkCopyMultiple/sampleSize=1024;numCopies=15-8      	  498631	     23668 ns/op	   64169 B/op	      73 allocs/op
BenchmarkCopyMultiple/sampleSize=1024;numCopies=20-8      	  384292	     30718 ns/op	   85626 B/op	      94 allocs/op
BenchmarkCopyMultiple/sampleSize=32768;numCopies=5-8      	  159147	     74279 ns/op	  656523 B/op	      32 allocs/op
BenchmarkCopyMultiple/sampleSize=32768;numCopies=10-8     	   86460	    137843 ns/op	 1312750 B/op	      53 allocs/op
BenchmarkCopyMultiple/sampleSize=32768;numCopies=15-8     	   60585	    197619 ns/op	 1968855 B/op	      73 allocs/op
BenchmarkCopyMultiple/sampleSize=32768;numCopies=20-8     	   46126	    257477 ns/op	 2625217 B/op	      94 allocs/op
BenchmarkCopyMultiple/sampleSize=1048576;numCopies=5-8    	    3901	   2829702 ns/op	20972689 B/op	      32 allocs/op
BenchmarkCopyMultiple/sampleSize=1048576;numCopies=10-8   	    2373	   5156780 ns/op	41945067 B/op	      53 allocs/op
BenchmarkCopyMultiple/sampleSize=1048576;numCopies=15-8   	    1681	   7565189 ns/op	62917331 B/op	      73 allocs/op
BenchmarkCopyMultiple/sampleSize=1048576;numCopies=20-8   	    1300	   9269650 ns/op	83889844 B/op	      94 allocs/op
BenchmarkIOCopy/sampleSize=32;numCopies=5-8               	15833799	       750.6 ns/op	    1120 B/op	      15 allocs/op
BenchmarkIOCopy/sampleSize=32;numCopies=10-8              	 8038702	      1497 ns/op	    2240 B/op	      30 allocs/op
BenchmarkIOCopy/sampleSize=32;numCopies=15-8              	 5369326	      2252 ns/op	    3360 B/op	      45 allocs/op
BenchmarkIOCopy/sampleSize=32;numCopies=20-8              	 3999177	      2981 ns/op	    4480 B/op	      60 allocs/op
BenchmarkIOCopy/sampleSize=1024;numCopies=5-8             	 3094920	      3869 ns/op	   20960 B/op	      15 allocs/op
BenchmarkIOCopy/sampleSize=1024;numCopies=10-8            	 1561958	      7683 ns/op	   41920 B/op	      30 allocs/op
BenchmarkIOCopy/sampleSize=1024;numCopies=15-8            	 1000000	     11577 ns/op	   62880 B/op	      45 allocs/op
BenchmarkIOCopy/sampleSize=1024;numCopies=20-8            	  765771	     15431 ns/op	   83840 B/op	      60 allocs/op
BenchmarkIOCopy/sampleSize=32768;numCopies=5-8            	  140972	     84595 ns/op	  655842 B/op	      15 allocs/op
BenchmarkIOCopy/sampleSize=32768;numCopies=10-8           	   70680	    168910 ns/op	 1311685 B/op	      30 allocs/op
BenchmarkIOCopy/sampleSize=32768;numCopies=15-8           	   47150	    253637 ns/op	 1967528 B/op	      45 allocs/op
BenchmarkIOCopy/sampleSize=32768;numCopies=20-8           	   35425	    339017 ns/op	 2623371 B/op	      60 allocs/op
BenchmarkIOCopy/sampleSize=1048576;numCopies=5-8          	    2545	   4765319 ns/op	20972025 B/op	      15 allocs/op
BenchmarkIOCopy/sampleSize=1048576;numCopies=10-8         	    1338	   9500500 ns/op	41944051 B/op	      30 allocs/op
BenchmarkIOCopy/sampleSize=1048576;numCopies=15-8         	     879	  13613407 ns/op	62916084 B/op	      45 allocs/op
BenchmarkIOCopy/sampleSize=1048576;numCopies=20-8         	     681	  18310750 ns/op	83888107 B/op	      61 allocs/op
```
