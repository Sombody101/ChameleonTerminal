package main

import (
	"fmt"
	"os"

	"gecko/argparse"
	"gecko/colors"
)

func main() {
	args := os.Args[1:]

	// Sort text and command switches, then get a config
	conf := argparse.ParseArguments(args)

	// Get a list of supported colors (unless otherwise specified by user)
	pcolors := colors.LoadPossibleColors(os.Getenv("TERM"), conf)

	end := ""
	if conf.NewLine {
		end = "\n"
	}

	var markupText string
	if conf.NoMarkup {
		markupText = conf.TextInput
	} else {
		markupText = colors.MarkupText(pcolors, conf)
	}

	fmt.Printf("%s%s", markupText, end)
}
