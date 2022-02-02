package main

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

type AddCommand struct {
	All bool `short:"a" long:"all" description:"Add all files"`
}

func (x *AddCommand) Execute(args []string) error {
	fmt.Printf("Adding (all=%v): %#v\n", x.All, args)
	return nil
}

func initAdd(parser *flags.Parser) {
	addCommand := new(AddCommand)
	parser.AddCommand("add",
		"Add a file",
		"The add command adds a file to the repository. Use -a to add all files.",
		addCommand)
}
