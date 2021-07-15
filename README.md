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
BenchmarkCopyMultiple/sampleSize=32;numCopies=5-8         	  508651	     23020 ns/op	  165817 B/op	      48 allocs/op
BenchmarkCopyMultiple/sampleSize=32;numCopies=10-8        	  288447	     41557 ns/op	  331228 B/op	      84 allocs/op
BenchmarkCopyMultiple/sampleSize=32;numCopies=15-8        	  195666	     61675 ns/op	  496513 B/op	     119 allocs/op
BenchmarkCopyMultiple/sampleSize=32;numCopies=20-8        	  147198	     80087 ns/op	  662054 B/op	     155 allocs/op
BenchmarkCopyMultiple/sampleSize=1024;numCopies=5-8       	  388245	     29584 ns/op	  185657 B/op	      48 allocs/op
BenchmarkCopyMultiple/sampleSize=1024;numCopies=10-8      	  231361	     52938 ns/op	  370909 B/op	      84 allocs/op
BenchmarkCopyMultiple/sampleSize=1024;numCopies=15-8      	  165648	     73975 ns/op	  556036 B/op	     119 allocs/op
BenchmarkCopyMultiple/sampleSize=1024;numCopies=20-8      	  118911	     98788 ns/op	  741419 B/op	     155 allocs/op
BenchmarkCopyMultiple/sampleSize=32768;numCopies=5-8      	   56552	    215397 ns/op	 1967432 B/op	      58 allocs/op
BenchmarkCopyMultiple/sampleSize=32768;numCopies=10-8     	   30182	    404303 ns/op	 3934466 B/op	     104 allocs/op
BenchmarkCopyMultiple/sampleSize=32768;numCopies=15-8     	   20986	    575410 ns/op	 5901409 B/op	     150 allocs/op
BenchmarkCopyMultiple/sampleSize=32768;numCopies=20-8     	   15963	    758097 ns/op	 7868580 B/op	     196 allocs/op
BenchmarkCopyMultiple/sampleSize=1048576;numCopies=5-8    	    1322	   8765747 ns/op	82412896 B/op	      83 allocs/op
BenchmarkCopyMultiple/sampleSize=1048576;numCopies=10-8   	     844	  14620172 ns/op	164825370 B/op	     154 allocs/op
BenchmarkCopyMultiple/sampleSize=1048576;numCopies=15-8   	     572	  21344404 ns/op	247237711 B/op	     224 allocs/op
BenchmarkCopyMultiple/sampleSize=1048576;numCopies=20-8   	     439	  26737609 ns/op	329650294 B/op	     296 allocs/op
BenchmarkIOCopy/sampleSize=32;numCopies=5-8               	 8845692	      1270 ns/op	    1200 B/op	      25 allocs/op
BenchmarkIOCopy/sampleSize=32;numCopies=10-8              	 4844488	      2473 ns/op	    2400 B/op	      50 allocs/op
BenchmarkIOCopy/sampleSize=32;numCopies=15-8              	 3240243	      3684 ns/op	    3600 B/op	      75 allocs/op
BenchmarkIOCopy/sampleSize=32;numCopies=20-8              	 2452618	      4900 ns/op	    4800 B/op	     100 allocs/op
BenchmarkIOCopy/sampleSize=1024;numCopies=5-8             	 2206700	      5415 ns/op	   21040 B/op	      25 allocs/op
BenchmarkIOCopy/sampleSize=1024;numCopies=10-8            	 1000000	     10879 ns/op	   42080 B/op	      50 allocs/op
BenchmarkIOCopy/sampleSize=1024;numCopies=15-8            	  718512	     16222 ns/op	   63120 B/op	      75 allocs/op
BenchmarkIOCopy/sampleSize=1024;numCopies=20-8            	  549104	     22375 ns/op	   84160 B/op	     100 allocs/op
BenchmarkIOCopy/sampleSize=32768;numCopies=5-8            	   46941	    254475 ns/op	 1802809 B/op	      35 allocs/op
BenchmarkIOCopy/sampleSize=32768;numCopies=10-8           	   24010	    502865 ns/op	 3605617 B/op	      70 allocs/op
BenchmarkIOCopy/sampleSize=32768;numCopies=15-8           	   15583	    766698 ns/op	 5408424 B/op	     105 allocs/op
BenchmarkIOCopy/sampleSize=32768;numCopies=20-8           	   10000	   1004489 ns/op	 7211234 B/op	     140 allocs/op
BenchmarkIOCopy/sampleSize=1048576;numCopies=5-8          	     871	  13611072 ns/op	82248330 B/op	      61 allocs/op
BenchmarkIOCopy/sampleSize=1048576;numCopies=10-8         	     440	  26942669 ns/op	164496664 B/op	     122 allocs/op
BenchmarkIOCopy/sampleSize=1048576;numCopies=15-8         	     296	  40113596 ns/op	246745017 B/op	     183 allocs/op
BenchmarkIOCopy/sampleSize=1048576;numCopies=20-8         	     218	  53434930 ns/op	328993374 B/op	     244 allocs/op
```
