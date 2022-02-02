package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/geomyidia/flagswrap"
	"github.com/jessevdk/go-flags"
)

type EditorOptions struct {
	Input  flags.Filename `short:"i" long:"input" description:"Input file" default:"-"`
	Output flags.Filename `short:"o" long:"output" description:"Output file" default:"-"`
}

type Point struct {
	X, Y int
}

func (p *Point) UnmarshalFlag(value string) error {
	parts := strings.Split(value, ",")

	if len(parts) != 2 {
		return errors.New("expected two numbers separated by a ,")
	}

	x, err := strconv.ParseInt(parts[0], 10, 32)

	if err != nil {
		return err
	}

	y, err := strconv.ParseInt(parts[1], 10, 32)

	if err != nil {
		return err
	}

	p.X = int(x)
	p.Y = int(y)

	return nil
}

func (p Point) MarshalFlag() (string, error) {
	return fmt.Sprintf("%d,%d", p.X, p.Y), nil
}

type Options struct {
	// Example of verbosity with level
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`

	// Example of optional value
	User string `short:"u" long:"user" description:"User name" optional:"yes" optional-value:"pancake"`

	// Example of map with multiple default values
	Users map[string]string `long:"users" description:"User e-mail map" default:"system:system@example.org" default:"admin:admin@example.org"`

	// Example of option group
	Editor EditorOptions `group:"Editor Options"`

	// Example of custom type Marshal/Unmarshal
	Point Point `long:"point" description:"A x,y point" default:"1,2"`

	Version bool `long:"version" description:"Get the current version of flagswrap"`
}

func main() {
	var options Options
	var parser = flags.NewParser(&options, flags.Default)
	parser.SubcommandsOptional = true

	initAdd(parser)
	initRm(parser)
	args, err := parser.Parse()
	if err != nil {
		wrappedErr := flagswrap.WrapError(err)
		switch {
		case wrappedErr.IsHelp():
			os.Exit(0)
		// Self-documenting go-flags errors:
		case wrappedErr.IsVerbose():
			os.Exit(1)
		// go-flags errors that need more context:
		case wrappedErr.IsSilent():
			fmt.Printf("Error: %v\n", wrappedErr)
			os.Exit(1)
		default:
			fmt.Printf("Error (unexpected): %+v\n", wrappedErr)
			os.Exit(1)
		}
	}
	if parser.FindOptionByLongName("version").Value().(bool) {
		fmt.Printf("%s\n", flagswrap.Version())
		os.Exit(0)
	}
	if parser.Active == nil {
		if len(args) > 0 {
			fmt.Printf("%v: %s\n", flags.ErrUnknownCommand, args[0])
		} else {
			fmt.Printf("%v\n", flags.ErrCommandRequired)
		}
		os.Exit(1)
	}
	fmt.Printf("%v\n", flags.ErrUnknown)
	os.Exit(1)
}