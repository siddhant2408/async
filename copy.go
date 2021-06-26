// Package async wraps golang library methods and provides a way to call them in a non blocking way.
package async

import (
	"context"
	"errors"
	"io"
)

const bufferSize = 1 << 15

// ErrorHandle must be provided by the caller as a function to handle the error coming from the goroutine should it occur.
type ErrorHandle func(err error)

// Copy provides a non blocking way to copy data from reader to writer.
func Copy(ctx context.Context, dst io.Writer, src io.Reader, errHandle ErrorHandle) func() {
	return CopyBuffer(ctx, dst, src, nil, errHandle)
}

// CopyBuffer provides a non blocking way to copy data from reader to writer using a provided buffer
// The callback function returned lets you handle the error by calling the ErrorHandle function provided during the call.
// This provides a freedom to handle the error whenever we want to.
// Please note that the returned function is a blocking one.
func CopyBuffer(ctx context.Context, dst io.Writer, src io.Reader, buf []byte, errHandle ErrorHandle) func() {
	if buf == nil {
		buf = make([]byte, bufferSize)
	}
	errChan := make(chan error, 1)
	go copy(ctx, dst, src, buf, errChan)
	return func() {
		err, _ := <-errChan
		if err != nil {
			errHandle(err)
		}
	}
}

func copy(ctx context.Context, dst io.Writer, src io.Reader, buf []byte, errChan chan error) {
	defer close(errChan)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := copyStream(dst, src, buf)
			if err != nil {
				if err != io.EOF {
					errChan <- err
				}
				return
			}
		}
	}
}

// copystream copies data from src to dst equal to the size of the buffer
func copyStream(dst io.Writer, src io.Reader, buf []byte) error {
	nr, er := src.Read(buf)
	if nr > 0 {
		nw, ew := dst.Write(buf[0:nr])
		if nw < 0 || nr < nw {
			nw = 0
			if ew == nil {
				return errors.New("invalid write result")
			}
		}
		if ew != nil {
			return ew
		}
		if nr != nw {
			return io.ErrShortWrite
		}
	}
	if er != nil {
		return er
	}
	return nil
}
