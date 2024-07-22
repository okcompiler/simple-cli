package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		args []string
		config
		output string
		err    error
	}{
		{
			args:   []string{"-h"},
			config: config{numTimes: 0},
			output: `
A greeter application which prints the name you entered a specified
number of times.

Usage of greeter: <options> [name]

Options:
  -n int
    	Number of times to greet
  -o string
    	Create an HTML document at the file path specified
`,
			err: errors.New("flag: help requested"),
		},
		{
			args:   []string{"-n", "10"},
			config: config{numTimes: 10},
			err:    nil,
		},
		{
			args:   []string{"-n", "abc"},
			config: config{numTimes: 0},
			err:    errors.New("invalid value \"abc\" for flag -n: parse error"),
		},
		{
			args:   []string{"-n", "1", "First Last"},
			config: config{numTimes: 1, name: "First Last"},
			err:    nil,
		},
		{
			args:   []string{"-n", "1", "First", "Last"},
			config: config{numTimes: 1},
			err:    errors.New("more than one positional argument specified"),
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range tests {
		c, err := parseArgs(byteBuf, tc.args)
		if tc.err == nil && err != nil {
			t.Fatalf("Expected nil error, got: %v\n", err)
		}
		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error to be: %v, got: %v\n", tc.err, err)
		}
		if c.numTimes != tc.numTimes {
			t.Errorf("Expected numTimes to be: %v, got: %v\n", tc.numTimes, c.numTimes)
		}
		gotMsg := byteBuf.String()
		if len(tc.output) != 0 && gotMsg != tc.output {
			t.Errorf("\nExpected stdout message to be:\n%#v,\ngot:\n%#v\n", tc.output, gotMsg)
		}
		byteBuf.Reset()
	}
}
