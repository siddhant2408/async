// Package async wraps golang library methods and provides a way to call them in a non blocking way.
package async

import (
	"context"
	"errors"
	"io"
)

const bufferSize = 1 << 15

// Copy provides a non blocking way to copy data from reader to writer.
func Copy(ctx context.Context, dst io.Writer, src io.Reader) ErrorHandle {
	return CopyBuffer(ctx, dst, src, nil)
}

// CopyBuffer provides a non blocking way to copy data from reader to writer using a provided buffer.
// The callback interface returned lets you handle the error async by calling the Err method.
// The Err method can be called immediately for synchronous behaviour or deferred to handle the error asynchronously.
func CopyBuffer(ctx context.Context, dst io.Writer, src io.Reader, buf []byte) ErrorHandle {
	if buf == nil {
		buf = make([]byte, bufferSize)
	}
	errChan := make(errorChannel, 1)
	go copy(ctx, dst, src, buf, errChan)
	return errChan
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
