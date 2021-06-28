package async

import (
	"context"
	"io"
	"sync"
)

const bufferSize = 1 << 15

// TeeReader represents a pair of Reader and Writer to perform the copy.
// It's the job of the caller to validate the reader and writer interfaces.
type TeeReader struct {
	Writer io.Writer
	Reader io.Reader
}

// Copy provides a non blocking way to copy data from reader to writer.
func Copy(ctx context.Context, teeReader *TeeReader) ErrorHandle {
	return CopyWithBuffer(ctx, teeReader, nil)
}

// CopyWithBuffer provides a non blocking way to copy data from reader to writer using a provided buffer.
// The returned function can be called immediately for synchronous behaviour or deferred to handle the error asynchronously.
func CopyWithBuffer(ctx context.Context, teeReader *TeeReader, buf []byte) ErrorHandle {
	errChan := make(chan error, 1)
	go copySingle(ctx, teeReader, buf, errChan)
	return func() error {
		return <-errChan
	}
}

// CopyMultiple performs concurrent copy of the tuples.
// It returns the first error that occurs from the concurrent copy.
// All copy goroutines are terminated as soon as an error occurs from any one of them.
// Use Copy or CopyBuffer functions for single copy as this starts additional goroutines.
func CopyMultiple(ctx context.Context, teeReaders []*TeeReader) ErrorHandle {
	ctx, cancel := context.WithCancel(ctx)
	mainErrChan := make(chan error, 1)
	errChan := make(chan error, len(teeReaders))
	wg := new(sync.WaitGroup)
	for _, reader := range teeReaders {
		wg.Add(1)
		go copyMultiple(ctx, reader, errChan, wg)
	}
	go runCopyManager(ctx, cancel, errChan, mainErrChan, wg)
	return func() error {
		return <-mainErrChan
	}
}

func copyMultiple(ctx context.Context, teeReader *TeeReader, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	copy(ctx, teeReader, nil, errChan)
}

// copy manager waits till all copy goroutine operations are done and then closes the main error channel
func runCopyManager(ctx context.Context, cancel context.CancelFunc, errChan chan error, mainErrChan chan error, wg *sync.WaitGroup) {
	defer close(mainErrChan)
	go runErrorListener(errChan, mainErrChan, cancel)
	wg.Wait()
	close(errChan)
}

// close all goroutines if an error comes from any goroutine performing the copy
func runErrorListener(errChan chan error, mainErrChan chan error, cancel context.CancelFunc) {
	err := <-errChan
	if err != nil {
		mainErrChan <- err
		cancel()
	}
}

func copySingle(ctx context.Context, teeReader *TeeReader, buf []byte, errChan chan error) {
	defer close(errChan)
	copy(ctx, teeReader, buf, errChan)
}

func copy(ctx context.Context, teeReader *TeeReader, buf []byte, errChan chan error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	if wt, ok := teeReader.Reader.(io.WriterTo); ok {
		_, err := wt.WriteTo(teeReader.Writer)
		if err != nil {
			errChan <- err
		}
		return
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := teeReader.Writer.(io.ReaderFrom); ok {
		_, err := rt.ReadFrom(teeReader.Reader)
		if err != nil {
			errChan <- err
		}
		return
	}
	if buf == nil {
		buf = make([]byte, bufferSize)
	}
	reader := io.TeeReader(teeReader.Reader, teeReader.Writer)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, err := reader.Read(buf)
			if err != nil {
				if err != io.EOF {
					errChan <- err
				}
				return
			}
		}
	}
}
