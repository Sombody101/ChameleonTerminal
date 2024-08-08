package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"gecko/argparse"
	"gecko/colors"
)

func main() {
	// Remove app path
	args := os.Args[1:]

	// Sort text and command switches, then get a config
	argparse.ParseArguments(args)

	// Get a list of supported colors (unless otherwise specified by user)
	colors.LoadPossibleColors(os.Getenv("TERM"))

	if argparse.Configuration.PrintHelp {
		colors.PrintHelpInfo()
		os.Exit(0)
	}

	// List all colors
	if argparse.Configuration.ListColors {
		listColors()
		os.Exit(0)
	}

	// List all styles
	if argparse.Configuration.ListStyles {
		listStyles()
		os.Exit(0)
	}

	// Add newline if selected
	end := '\000'
	if argparse.Configuration.NewLine {
		end = '\n'
	}

	// Markup line if selected
	var markupText string
	if argparse.Configuration.NoMarkup {
		markupText = argparse.Configuration.TextInput
	} else {
		markupText = colors.MarkupText(argparse.Configuration.TextInput)
	}

	// Print final output
	fmt.Printf("%s%c", markupText, end)
}

// List all supported colors
func listColors() {

	// Sample text to be displayed [Overridable]
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
		fmt.Println(colors.MarkupText(fmt.Sprintf(
			"%-17s (#%06x) [%s]%s[/]",
			colorName,
			colorCode,
			sampleColor,
			sampleText)))
	}
}

func listStyles() {

	// Sample text to be displayed [Overridable]
	sampleText := "Hello, World!"
	if argparse.Configuration.TextInput != "" {
		sampleText = argparse.Configuration.TextInput
	}

	// Load all colors into an array so they can be sorted
	keys := make([]string, 0, len(colors.FormatCodes))
	for k := range colors.FormatCodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Print all colors to the console with the sample text
	for _, styleName := range keys {
		sampleStyle := styleName

		styleCode := strings.Trim(colors.FormatCodes[styleName], ";")
		fmt.Println(colors.MarkupText(fmt.Sprintf(
			"%-11s (ANSI: %+2s) [%s]%s[/]",
			styleName,
			styleCode,
			sampleStyle,
			sampleText)))
	}
}
