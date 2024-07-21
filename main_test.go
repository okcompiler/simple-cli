package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"testing"
	"time"
)

var binaryName string

func TestMain(m *testing.M) {
	if runtime.GOOS == "windows" {
		binaryName = "test-binary.exe"
	} else {
		binaryName = "test-binary"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "build", "-o", binaryName)
	err := cmd.Run()
	if err != nil {
		os.Exit(1)
	}
	defer func() {
		err = os.Remove(binaryName)
		if err != nil {
			log.Fatalf("Error removing built binary: %v", err)
		}
	}()

	m.Run()
}

func TestApplication(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	defer cancel()

	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	binaryPath := path.Join(currentDir, binaryName)

	tests := []struct {
		args                []string
		input               string
		expectedOutputLines []string
		expectedExitCode    int
	}{
		{
			args: []string{},
			expectedOutputLines: []string{
				"must specify a number greater than 0",
			},
			expectedExitCode: 1,
		},
		{
			args: []string{"-h"},
			expectedOutputLines: []string{
				"flag: help requested",
			},
			expectedExitCode: 1,
		},
		{
			args: []string{"a"},
			expectedOutputLines: []string{
				"positional arguments specified",
			},
			expectedExitCode: 1,
		},
		{
			args:  []string{"-n", "2"},
			input: "First Last",
			expectedOutputLines: []string{
				"Your name please? Press the Enter key when done.",
				"Nice to meet you First Last",
				"Nice to meet you First Last",
			},
			expectedExitCode: 0,
		},
	}

	byteBuf := new(bytes.Buffer)
	for _, tc := range tests {
		cmd := exec.CommandContext(ctx, binaryPath, tc.args...)
		cmd.Stdout = byteBuf
		if len(tc.input) != 0 {
			cmd.Stdin = strings.NewReader(tc.input)
		}

		err := cmd.Run()
		if err != nil && tc.expectedExitCode == 0 {
			t.Fatalf("Expected application to exit without an error. Got: %v", err)
		}
		if cmd.ProcessState.ExitCode() != tc.expectedExitCode {
			t.Log(byteBuf.String())
			t.Fatalf(
				"Expected application to have exit code: %v. Got: %v",
				tc.expectedExitCode,
				cmd.ProcessState.ExitCode(),
			)
		}

		output := byteBuf.String()
		lines := strings.Split(output, "\n")
		for num := range tc.expectedOutputLines {
			if lines[num] != tc.expectedOutputLines[num] {
				t.Fatalf(
					"Expected output line to be: %v, Got: %v",
					tc.expectedOutputLines[num],
					lines[num],
				)
			}
		}

		byteBuf.Reset()
	}
}
