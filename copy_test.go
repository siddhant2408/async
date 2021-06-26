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
	errHandle := func(err error) {}
	errCallback := async.Copy(context.Background(), dst, src, async.ErrorHandle(errHandle))
	errCallback()
	actual := dst.String()
	if expected != actual {
		t.Fatalf("unexpected result, got %s want %s", actual, expected)
	}
}

func TestCopyBuffer(t *testing.T) {
	expected := "test"
	src := bytes.NewReader([]byte(expected))
	dst := new(bytes.Buffer)
	errHandle := func(err error) {}
	buf := make([]byte, 1<<8)
	errCallback := async.CopyBuffer(context.Background(), dst, src, buf, async.ErrorHandle(errHandle))
	errCallback()
	actual := dst.String()
	if expected != actual {
		t.Fatalf("unexpected result, got %s want %s", actual, expected)
	}
}
