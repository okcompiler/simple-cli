package main

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

type errReader struct {
	err error
}

func (e errReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}

func TestGetName(t *testing.T) {
	tests := []struct {
		name string
		r    io.Reader
	}{
		{
			name: "First Last",
			r:    strings.NewReader("First Last"),
		},
		{
			r: errReader{err: errors.New("simulated reader error")},
		},
	}

	w := new(bytes.Buffer)
	for _, tc := range tests {
		n, err := getName(tc.r, w)
		if err != nil && err.Error() != "simulated reader error" {
			t.Fatalf("Expected simulated reader error, got: %v\n", err.Error())
		}
		if n != tc.name {
			t.Errorf("Expected name %v, got: %v", tc.name, n)
		}
	}
}
