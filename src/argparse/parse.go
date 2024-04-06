package argparse

import (
	"fmt"
	"gecko/version"
	"os"
	"strings"
)

type Config struct {
	// 0 = default,
	// 1 = force no color,
	// 2 = force color (256 colors),
	// 3 = force color (8 colors)
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

	// For debugging
	Verbose bool
}

var Configuration = Config{
	ForceColor:             0,
	Escape:                 false,
	NewLine:                true,
	NoMarkup:               false,
	TextInput:              "",
	PrintHelp:              false,
	ListColors:             false,
	ListColorsAsBackground: false,
	Verbose:                false,
}

/*
 * 	-n: No newline
 * 	-e: Escape codes
 *	-c: Force color
 *  -C: Force no color
 *	-M: No markup
 */
func ParseArguments(args []string) {
	if len(args) == 0 {
		Configuration.TextInput = ""
		return
	}

	var sb strings.Builder

	for i, word := range args {

		if word == "--" {
			// End of input arguments
			sb.WriteString(strings.Join(args[i:], " "))
			Configuration.TextInput = sb.String()
			return
		}

		if strings.HasPrefix(word, "--") {
			// Add long-hand options
			switch word[2:] {
			case "help":
				Configuration.PrintHelp = true
				return

			case "listc":
				Configuration.ListColors = true
				continue

			case "listcb":
				Configuration.ListColors = true
				Configuration.ListColorsAsBackground = true
				continue

			case "lists":
				Configuration.ListStyles = true
				continue

			case "verbose":
				Configuration.Verbose = true
				continue

			case "version":
				fmt.Printf("v%s\nch: %s\nbt: %s\n", version.VERSION, version.COMMIT_HASH, version.BUILD_TIME)
				os.Exit(0)
			}

		} else if strings.HasPrefix(word, "-") {
			// Add short-hand options
			for _, c := range word[1:] {
				switch c {

				case 'n':
					Configuration.NewLine = false
					continue

				case 'e':
					Configuration.Escape = true
					continue

				case 'c':
					Configuration.ForceColor = 2 // 256 color
					continue

				//case 'z':
				//	Configuration.ForceColor = 3 // 8 color
				//	continue

				case 'C':
					Configuration.ForceColor = 1 // no color
					continue

				case 'M':
					Configuration.NoMarkup = true
					continue

				default:
					// Just append it to the output (like echo)

					//fmt.Printf("gecko: Unknown option '%c'\n", c)
					//os.Exit(2)
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
