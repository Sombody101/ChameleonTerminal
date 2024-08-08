package argparse

import (
	"fmt"
	"gecko/verbose"
	"gecko/version"
	"os"
	"strings"
)

type Config struct {
	// 0 = default,
	// 1 = force no color,
	// 2 = force color (256 colors),
	ForceColor byte

	// Resolve escape codes
	Escape bool

	// Add a newline at the end of args
	NewLine bool

	// Resolve markups
	NoMarkup bool

	// Input minus the config switches
	TextInput string

	// Print help info then exit
	PrintHelp bool

	ListColors bool

	ListColorsAsBackground bool

	ListStyles bool
}

var Configuration Config

/*
 * 	-n: No newline
 * 	-e: Escape codes
 *	-c: Force color
 *  -C: Force no color
 *	-M: No markup
 */
func ParseArguments(args []string) {

	// Initialize the configuration
	Configuration = Config{
		ForceColor:             0,
		Escape:                 false,
		NewLine:                true,
		NoMarkup:               false,
		TextInput:              "",
		PrintHelp:              false,
		ListColors:             false,
		ListColorsAsBackground: false,
	}

	if len(args) == 0 {
		Configuration.TextInput = ""
		return
	}

	// Create an output buffer for non-argument text
	var sb strings.Builder

	// Run through all text looking for switches
	for i, word := range args {

		if word == "--" {
			// Stop looking for switch arguments
			sb.WriteString(strings.Join(args[i:], " "))
			Configuration.TextInput = sb.String()
			return
		}

		// Add long-hand options
		if strings.HasPrefix(word, "--") {

			// Check if switch is known, and apply it to the configuration if so
			switch resolveSwitch(word) {
			case 0:
				continue // Don't add the switch to the output buffer

			case 1:
				return // Used for switches like 'help'

			case 3:
				// Treat this as input text that looked like a switch
			}

		} else /* Add short-hand options */ if strings.HasPrefix(word, "-") {
			for _, c := range word[1:] {
				// Check if switch is known, and apply it to the configuration if so
				switch resolveShorthandSwitch(c) {
				case 0:
					continue // Don't add the switch to the output buffer

				case 1:
					return // Used for switches like 'help'

				case 3:
					// Treat this as input text that looked like a switch
				}
			}

			continue
		}

		sb.WriteString(word)
		sb.WriteRune(' ')
	}

	output := sb.String()
	if len(output) != 0 {
		Configuration.TextInput = output[:len(output)-1]
	}
}

func resolveShorthandSwitch(char rune) byte {
	switch char {

	case 'n':
		Configuration.NewLine = false
		return 0

	case 'e':
		Configuration.Escape = true
		return 0

	case 'c':
		Configuration.ForceColor = 2 // 256 color
		return 0

	case 'C':
		Configuration.ForceColor = 1 // no color
		return 0

	case 'M':
		Configuration.NoMarkup = true
		return 0

	}

	verbose.Log("Unknown switch (treating as text)", string(char))
	return 3
}

// Returns: true when the given switch should be executed without parsing
// the rest of the arguments
func resolveSwitch(word string) byte {
	switch word[2:] {

	case "help":
		Configuration.PrintHelp = true
		return 1

	case "listc":
		Configuration.ListColors = true
		return 0

	case "listcb":
		Configuration.ListColors = true
		Configuration.ListColorsAsBackground = true
		return 0

	case "lists":
		Configuration.ListStyles = true
		return 0

	case "verbose":
		verbose.VerboseLoggingEnabled = true
		return 0

	case "version":
		fmt.Printf("v%s\nch: %s\nbt: %s\n", version.VERSION, version.COMMIT_HASH, version.BUILD_TIME)
		os.Exit(0)
		return 1
	}

	// Ignore unknown switches as text input, but log it for verbose mode
	// in case someone gets confused when they have a typo
	verbose.Log("Unknown switch (treating as text)", word)
	return 3
}
