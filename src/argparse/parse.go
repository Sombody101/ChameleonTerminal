package argparse

import (
	"fmt"
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
}

/*
 * 	-n: No newline
 * 	-e: Escape codes
 *	-c: Force color
 *  -C: Force no color
 *	-M: No markup
 */
func ParseArguments(args []string) Config {
	conf := Config{
		ForceColor: 0,
		Escape:     false,
		NewLine:    true,
		NoMarkup:   false,
	}

	var sb strings.Builder

	for i, word := range args {

		if word == "--" {
			// End of input arguments
			sb.WriteString(strings.Join(args[i:], " "))
			conf.TextInput = sb.String()
			return conf
		}

		if strings.HasPrefix(word, "--") {
			// Add long-hand options
			switch word[2:] {
			case "help":
				printHelpInfo()
				os.Exit(0)
			}

		} else if strings.HasPrefix(word, "-") {
			// Add short-hand options
			for _, c := range word[1:] {
				switch c {

				case 'n':
					conf.NewLine = false
					continue

				case 'e':
					conf.Escape = true
					continue

				case 'c':
					conf.ForceColor = 2 // 256 color
					continue

				case 'z':
					conf.ForceColor = 3 // 8 color
					continue

				case 'C':
					conf.ForceColor = 1 // no color
					continue

				case 'M':
					conf.NoMarkup = true
					continue

				default:
					fmt.Printf("gecko: Unknown option '%c'\n", c)
					os.Exit(2)
				}
			}

		}

		sb.WriteString(word)
		sb.WriteRune(' ')
	}

	output := sb.String()
	conf.TextInput = output[:len(output)-1]
	return conf
}

func printHelpInfo() {
	fmt.Println("usage: gecko [options...] [text..]")
	fmt.Println("Display text with markup (color), speedily and efficiently.")
	fmt.Println()
	fmt.Println("Stop reading options after standalone '--' is read.")
	fmt.Println()
	fmt.Println("    Options:")
	fmt.Println("      -n:\tNo newline when printing output text")
	fmt.Println()
	fmt.Println("      -e:\tResolve escape codes")
	fmt.Println()
	fmt.Println("      -c:\tForce color output even when not detected as supported (256 color)")
	fmt.Println()
	fmt.Println("      -z:\tForce color output even when not detected as supported (8 color)")
	fmt.Println()
	fmt.Println("      -C:\tForce no color output even when detected as supported")
	fmt.Println()
	fmt.Println("      -M:\tDo not resolve markup sequences")
	fmt.Println()
	fmt.Println("    Exit Status:")
	fmt.Println("\t2: Option not found")
	fmt.Println("\tdefault (0): Program success")

}
