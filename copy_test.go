package async_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/siddhant2408/async"
)

func TestCopy(t *testing.T) {
	expected := "test"
	src := &mockReader{bytes.NewReader([]byte(expected))}
	dst := &mockWriter{new(bytes.Buffer)}
	asyncErr := async.Copy(context.Background(), &async.TeeReader{dst, src})
	err := asyncErr()
	if err != nil {
		t.Fatal(err)
	}
	actual := dst.buf.String()
	if expected != actual {
		t.Fatalf("unexpected result, got %s want %s", actual, expected)
	}
}

func TestCopyBuffer(t *testing.T) {
	expected := "test"
	src := bytes.NewReader([]byte(expected))
	dst := &mockWriter{new(bytes.Buffer)}
	buf := make([]byte, 1<<8)
	asyncErr := async.CopyWithBuffer(context.Background(), &async.TeeReader{dst, src}, buf)
	err := asyncErr()
	if err != nil {
		t.Fatal(err)
	}
	actual := dst.buf.String()
	if expected != actual {
		t.Fatalf("unexpected result, got %s want %s", actual, expected)
	}
}

func TestCopyMultiple(t *testing.T) {
	expected1, expected2, expected3 := "test1", "test2", "test3"
	reader1, reader2, reader3 := bytes.NewReader([]byte(expected1)), bytes.NewReader([]byte(expected2)), bytes.NewReader([]byte(expected3))
	writer1, writer2, writer3 := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	asyncErr := async.CopyMultiple(context.Background(), []*async.TeeReader{
		{writer1, reader1},
		{writer2, reader2},
		{writer3, reader3},
	})
	err := asyncErr()
	if err != nil {
		t.Fatal(err)
	}
	actual1, actual2, actual3 := writer1.String(), writer2.String(), writer3.String()
	if expected1 != actual1 {
		t.Fatalf("unexpected result, got %s want %s", actual1, expected1)
	}
	if expected2 != actual2 {
		t.Fatalf("unexpected result, got %s want %s", actual2, expected2)
	}
	if expected3 != actual3 {
		t.Fatalf("unexpected result, got %s want %s", actual3, expected3)
	}
}

func BenchmarkCopyMultiple(b *testing.B) {
	sampleSize := []int{1 << 5, 1 << 10, 1 << 15, 1 << 20}
	numCopies := []int{5, 10, 15, 20}
	for _, curSampleSize := range sampleSize {
		for _, curNumCopy := range numCopies {
			b.Run(fmt.Sprintf("sampleSize=%d;numCopies=%d", curSampleSize, curNumCopy), func(b *testing.B) {
				runCopyMultiple(b, curNumCopy, curSampleSize)
			})
		}
	}
}

func runCopyMultiple(b *testing.B, numCopy int, sampleSize int) {
	s := new(strings.Builder)
	for i := 0; i < sampleSize; i++ {
		_, err := s.WriteString("test")
		if err != nil {
			b.Fatal(err)
		}
	}
	sbyte := []byte(s.String())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var tr []*async.TeeReader
		for i := 0; i < numCopy; i++ {
			tr = append(tr, &async.TeeReader{
				&mockWriter{new(bytes.Buffer)},
				&mockReader{bytes.NewReader(sbyte)}})
		}
		asyncErr := async.CopyMultiple(context.Background(), tr)
		err := asyncErr()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkIOCopy(b *testing.B) {
	sampleSize := []int{1 << 5, 1 << 10, 1 << 15, 1 << 20}
	numCopies := []int{5, 10, 15, 20}
	for _, curSampleSize := range sampleSize {
		for _, curNumCopy := range numCopies {
			b.Run(fmt.Sprintf("sampleSize=%d;numCopies=%d", curSampleSize, curNumCopy),
				func(b *testing.B) {
					runIOCopy(b, curNumCopy, curSampleSize)
				})
		}
	}
}

func runIOCopy(b *testing.B, numCopy int, sampleSize int) {
	s := new(strings.Builder)
	for i := 0; i < sampleSize; i++ {
		_, err := s.WriteString("test")
		if err != nil {
			b.Fatal(err)
		}
	}
	sbyte := []byte(s.String())
	buf := make([]byte, 1<<15)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < numCopy; i++ {
			reader := &mockReader{bytes.NewReader(sbyte)}
			writer := &mockWriter{new(bytes.Buffer)}
			_, err := io.CopyBuffer(writer, reader, buf)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

// these types are provided as an input to the methods as it does not implement io.WriteTo so
// that buffer allocation can happen
type mockReader struct {
	buf *bytes.Reader
}

func (m *mockReader) Read(p []byte) (int, error) {
	return m.buf.Read(p)
}

type mockWriter struct {
	buf *bytes.Buffer
}

func (m *mockWriter) Write(p []byte) (int, error) {
	return m.buf.Write(p)
}
