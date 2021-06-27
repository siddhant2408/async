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
	asyncErr := async.Copy(context.Background(), async.TeeReader{dst, src})
	err := asyncErr()
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
	asyncErr := async.CopyWithBuffer(context.Background(), async.TeeReader{dst, src}, buf)
	err := asyncErr()
	if err != nil {
		t.Fatal(err)
	}
	actual := dst.String()
	if expected != actual {
		t.Fatalf("unexpected result, got %s want %s", actual, expected)
	}
}

func TestCopyMultiple(t *testing.T) {
	expected1, expected2, expected3 := "test1", "test2", "test3"
	reader1, reader2, reader3 := bytes.NewReader([]byte(expected1)), bytes.NewReader([]byte(expected2)), bytes.NewReader([]byte(expected3))
	writer1, writer2, writer3 := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	asyncErr := async.CopyMultiple(context.Background(), []async.TeeReader{
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
