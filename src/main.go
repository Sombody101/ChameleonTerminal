package main

import (
	"fmt"
	"os"
	"sort"

	"gecko/argparse"
	"gecko/colors"
)

func main() {
	args := os.Args[1:]

	// Sort text and command switches, then get a config
	argparse.ParseArguments(args)

	// Get a list of supported colors (unless otherwise specified by user)
	colors.LoadPossibleColors(os.Getenv("TERM"))

	if argparse.Configuration.PrintHelp {
		colors.PrintHelpInfo()
	}

	if argparse.Configuration.ListColors {
		listColors()
	}

	// Add newline if selected
	end := ""
	if argparse.Configuration.NewLine {
		end = "\n"
	}

	// Markup line if selected
	var markupText string
	if argparse.Configuration.NoMarkup {
		markupText = argparse.Configuration.TextInput
	} else {
		markupText = colors.MarkupText(argparse.Configuration.TextInput)
	}

	// Print final output
	fmt.Printf("%s%s", markupText, end)
}

// List all supported colors
func listColors() {
	sampleText := "Hello, World!"
	if argparse.Configuration.TextInput != "" {
		sampleText = argparse.Configuration.TextInput
	}

	keys := make([]string, 0, len(colors.Colors))
	for k := range colors.Colors {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, colorName := range keys {
		sampleColor := colorName

		// If used with --listcb, print colors as background colors
		if argparse.Configuration.ListColorsAsBackground {
			sampleColor = fmt.Sprintf("_ on %v", colorName)
		}

		colorCode := colors.Colors[colorName]
		fmt.Println(colors.MarkupText(fmt.Sprintf("%-17s (#%06x) [%s]%s[/]", colorName, colorCode, sampleColor, sampleText)))
	}

	os.Exit(0)
}
