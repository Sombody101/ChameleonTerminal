package colors

import (
	"bufio"
	"fmt"
	"gecko/argparse"
	"os"
	"os/exec"
	"strings"
)

func LoadPossibleColors(term string) {

	switch argparse.Configuration.ForceColor {
	case 1:
		return // No color

	case 2:
		LoadColors256() // 256 colors
		return
	}

	if ansiSupported() {
		LoadColors256()
	}

	// Otherwise load no colors and print all text normally
}

func ansiSupported() bool {

	// Check using tput if available
	if _, err := exec.LookPath("tput"); err == nil {
		out, err := exec.Command("tput", "colors").Output()
		if err == nil {
			var numColors int

			_, err := fmt.Sscanln(string(out), &numColors)

			if err == nil && numColors >= 8 {
				return true
			}
		}
	}

	// Direct console query (ensure CSI is defined)
	const csi = "\033["
	fmt.Print(csi + "c")
	reader := bufio.NewReader(os.Stdin)
	ansiReport, err := reader.ReadString('c')

	if err != nil {
		return false
	}

	return strings.TrimSpace(ansiReport) != ""
}

func PrintHelpInfo() {
	fmt.Println("usage: gecko [[options...]] [[text..]]")
	fmt.Println("Display text with markup (color), speedily and efficiently.")
	fmt.Println()
	fmt.Println("Stops parsing options after standalone '--' is read.")
	fmt.Println()
	fmt.Println("    Options:")
	fmt.Println("      -n:\tNo newline when printing output text")
	//fmt.Println()
	//fmt.Println("      [gold1]-e[/]:\tResolve escape codes")
	fmt.Println()
	fmt.Println("      -c:\tForce color output even when not detected as supported (256 color)")
	fmt.Println()
	//fmt.Println("      [gold1]-z[/]:\tForce color output even when not detected as supported (8 color)")
	//fmt.Println()
	fmt.Println("      -C:\tForce no color output even when detected as supported")
	fmt.Println()
	fmt.Println("      -M:\tDo not resolve markup sequences")
	fmt.Println()
	fmt.Println("      --listc:\tList possible colors (listc[olors])")
	fmt.Println("          (Add text after to change sample text, markup is still parsed)")
	fmt.Println()
	fmt.Println("      --listcb:\tList possible colors as backgrounds (listc[olor]b[ackgrounds])")
	fmt.Println()
	fmt.Println("      --lists:\t List possible styles and their ANSI codes")
	fmt.Println()
	//fmt.Println("    Exit Status:")
	//fmt.Println("\t2: Option not found")
	//fmt.Println("\t0: Program success (default)")
}
