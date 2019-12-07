package safetyrecevier

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"testing"
)

func TestLargeFile(t *testing.T) {
	b, err := ioutil.ReadFile("example-gopher.png")
	exampleReader := bytes.NewReader(b)
	if err != nil {
		t.Error(err)
	}
	safeReceive := New(30) // make buffer size 30 bytes
	defer safeReceive.CloseReceive()

	for {
		chunk := make([]byte, 3)
		n, err := exampleReader.Read(chunk)
		safeReceive.Append(chunk, n)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Error(err)
		}
	}
	if err := validFile(safeReceive, b); err != nil {
		t.Error(err)
	}
}

func validFile(safeReceive *SafeReceive, target []byte) error {
	b, err := ioutil.ReadAll(safeReceive)
	if err != nil {
		return err
	}
	defer safeReceive.CloseReader()

	return valid(target, b)
}

func valid(b1, b2 []byte) error {
	fileInvalidErr := errors.New("file valid failed")
	if len(b1) != len(b2) {
		return fileInvalidErr
	}
	for i, v := range b1 {
		if b2[i] != v {
			return fileInvalidErr
		}
	}
	return nil
}

func TestSmall(t *testing.T) {
	b, err := ioutil.ReadFile("example-gopher.png")
	exampleReader := bytes.NewReader(b)
	if err != nil {
		t.Error(err)
	}
	safeReceive := New(1000000) // make buffer size 1 Mb
	defer safeReceive.CloseReceive()

	for {
		chunk := make([]byte, 10000)
		n, err := exampleReader.Read(chunk)
		safeReceive.Append(chunk, n)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Error(err)
		}
	}
	b2, err := ioutil.ReadAll(safeReceive)
	if err != nil {
		t.Error(err)
	}
	if err := valid(b, b2); err != nil {
		t.Error(err)
	}
}
