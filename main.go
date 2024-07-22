package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
)

type config struct {
	numTimes     int
	htmlFilePath string
	name         string
}

var errInvalidPosArgSpecified = errors.New("more than one positional argument specified")

func parseArgs(w io.Writer, args []string) (config, error) {
	c := config{}
	fs := flag.NewFlagSet("greeter", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.Usage = func() {
		usageString := `
A greeter application which prints the name you entered a specified
number of times.

Usage of %s: <options> [name]

`
		fmt.Fprintf(w, usageString, fs.Name())
		fmt.Fprintln(w, "Options:")
		fs.PrintDefaults()
	}
	fs.IntVar(&c.numTimes, "n", 0, "Number of times to greet")
	fs.StringVar(&c.htmlFilePath, "o", "", "Create an HTML document at the file path specified")
	err := fs.Parse(args)
	if err != nil {
		return c, err
	}
	if fs.NArg() > 1 {
		return c, errInvalidPosArgSpecified
	}
	if fs.NArg() == 1 {
		c.name = fs.Arg(0)
	}
	return c, nil
}

func validateArgs(c config) error {
	if c.numTimes <= 0 && len(c.htmlFilePath) == 0 {
		return errors.New("must specify a number greater than 0")
	}

	return nil
}

func runCmd(r io.Reader, w io.Writer, c config) error {
	var err error
	if len(c.name) == 0 {
		c.name, err = getName(r, w)
		if err != nil {
			return err
		}
	}

	if len(c.htmlFilePath) != 0 {
		return greetWithHTML(c.htmlFilePath, c.name)
	}

	greetUser(c, w)
	return nil
}

func greetWithHTML(path, name string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl, err := template.New("greeterHTML").Parse("<h1>Hello {{.}}</h1>")
	if err != nil {
		return err
	}
	return tmpl.Execute(f, name)
}

func greetUser(c config, w io.Writer) {
	msg := fmt.Sprintf("Nice to meet you %s\n", c.name)

	for i := 0; i < c.numTimes; i++ {
		fmt.Fprint(w, msg)
	}
}

func getName(r io.Reader, w io.Writer) (string, error) {
	msg := "Your name please? Press the Enter key when done.\n"
	fmt.Fprint(w, msg)

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	name := scanner.Text()
	if len(name) == 0 {
		return "", errors.New("you didn't enter your name")
	}

	return name, nil
}

func main() {
	c, err := parseArgs(os.Stderr, os.Args[1:])
	if err != nil {
		if errors.Is(err, errInvalidPosArgSpecified) {
			fmt.Fprintln(os.Stdout, err)
		}
		os.Exit(1)
	}
	err = validateArgs(c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
	err = runCmd(os.Stdin, os.Stdout, c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
