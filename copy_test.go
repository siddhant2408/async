package async_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/siddhant2408/async"
)

func TestCopy(t *testing.T) {
	expected := "test"
	src := bytes.NewReader([]byte(expected))
	dst := new(bytes.Buffer)
	asyncErr := async.Copy(context.Background(), dst, src)
	err := asyncErr.Err()
	if err != nil {
		t.Fatal(err)
	}
	actual := dst.String()
	if expected != actual {
		t.Fatalf("unexpected result, got %s want %s", actual, expected)
	}
}

func TestCopyBuffer(t *testing.T) {
	expected := "test"
	src := bytes.NewReader([]byte(expected))
	dst := new(bytes.Buffer)
	buf := make([]byte, 1<<8)
	asyncErr := async.CopyBuffer(context.Background(), dst, src, buf)
	err := asyncErr.Err()
	if err != nil {
		t.Fatal(err)
	}
	actual := dst.String()
	if expected != actual {
		t.Fatalf("unexpected result, got %s want %s", actual, expected)
	}
}
