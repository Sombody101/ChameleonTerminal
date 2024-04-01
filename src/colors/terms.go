package colors

import (
	"fmt"
	"gecko/argparse"
	"os/exec"

	"strings"
)

func LoadPossibleColors(term string, conf argparse.Config) map[string]int {
	pcolors := map[string]int{}

	switch conf.ForceColor {
	case 1:
		return pcolors // No color

	case 2:
		LoadColors256(pcolors) // 256 colors
		return pcolors

	case 3:
		LoadColors8(pcolors) // 8 colors
		return pcolors
	}

	b_colors, err := exec.Command("tput", "colors").Output()

	if err == nil {
		s_colors := string(b_colors[:len(b_colors)-1]) // Trim newline
		//colors, s_err := strconv.Atoi(string(s_colors))
		//if s_err != nil {
		//	return manualColorCheck(term, pcolors)
		//}

		switch s_colors {
		case "256":
			LoadColors256(pcolors)

		case "8":
			LoadColors8(pcolors)

		default:
			// Try manually getting the colors
			manualColorCheck(term, pcolors)
		}

	} else {
		return manualColorCheck(term, pcolors)
	}

	return pcolors
}

func manualColorCheck(term string, colors map[string]int) map[string]int {

	//if strings.Contains(strings.ToLower(term), "256color") {
	//	LoadColors256(colors)
	//} else {
	//	LoadColors8(colors)
	//}

	fmt.Println("manual")

	/*
		xterm: Supports up to 256 colors (xterm-256color).
		rxvt: Supports up to 64 colors (rxvt-256color for 256 colors).
		rxvt-unicode: Supports up to 256 colors (rxvt-unicode-256color).
		linux: Typically supports only 16 colors.
		vt100 and vt220: Usually support only basic ANSI colors (8 colors).
		ansi: Supports basic ANSI colors (8 colors).
		screen: Supports up to 256 colors (screen-256color).
		tmux: Supports up to 256 colors (tmux-256color).
		konsole: Supports up to 256 colors (konsole-256color).
		gnome-256color: Supports up to 256 colors.
		terminator: Supports up to 256 colors.
	*/

	switch strings.ToLower(term) {

	// 256 colors
	case "xterm":
	case "rxvt-unicode":
	case "rxvt-unicode-256color":
		LoadColors256(colors)

	// 8 colors
	case "ansi":
	case "rxvt":
	case "linux":
	case "vs100":
	case "vt220":
		LoadColors8(colors)

	default:
		// Otherwise, load no colors
	}

	return colors
}
