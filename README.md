# Async

## Introduction

Async wraps golang library methods and provides a way to call them in a non blocking way. The context parameter adds more control over the lifecycle of the goroutine.

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
