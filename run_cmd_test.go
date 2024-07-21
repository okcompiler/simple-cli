package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		c      config
		input  string
		output string
		err    error
	}{
		{
			c:      config{numTimes: 5},
			input:  "",
			output: strings.Repeat("Your name please? Press the Enter key when done.\n", 1),
			err:    errors.New("you didn't enter your name"),
		},
		{
			c:     config{numTimes: 5},
			input: "First Last",
			output: "Your name please? Press the Enter key when done.\n" + strings.Repeat(
				"Nice to meet you First Last\n",
				5,
			),
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range tests {
		r := strings.NewReader(tc.input)
		err := runCmd(r, byteBuf, tc.c)
		if err != nil && tc.err == nil {
			t.Fatalf("Expected nil error, got: %v\n", err)
		}
		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error: %v, got error: %v\n", tc.err.Error(), err.Error())
		}
		gotMsg := byteBuf.String()
		if gotMsg != tc.output {
			t.Errorf("Exptected stdout message to be: %v, got: %v\n", tc.output, gotMsg)
		}

		byteBuf.Reset()
	}
}
