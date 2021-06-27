# Async

## Introduction

Async wraps golang library methods and provides a way to call them in a non blocking way.

## Basic Usage

```go
go get github.com/siddhant2408/async
```

## Copy

```go
errHandle := async.ErrorHandle(func(err error) {
	//handle the error
})
callback := async.Copy(ctx, dst, src, errHandle)
defer callback()
```

## Callback

The callback method lets you handle the error coming from the method, if any, in a non blocking way in either of the 2 ways:
1. You can call this method before you need the result as this method is blocking.
2. Defer it till the end of the caller lifecycle. This lets you handle the error in an asynchronous way.
