package colors

import (
	"fmt"
	"image/color"
	"math"
	"runtime"
	"strconv"
	"strings"

	"gecko/argparse"
)

const (
	RESET_COLOR = "\033[0m"
)

func IntHexToRGB(hexColor int) (byte, byte, byte) {
	r := byte((hexColor >> 16) & 0xFF)
	g := byte((hexColor >> 8) & 0xFF)
	b := byte(hexColor & 0xFF)
	return r, g, b
}

func StringHexToRGB(hexColor string) (int, int, int) {
	hex := hexColor[1:]
	r, g, b, _ := color.RGBA{}.RGBA()
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return int(r / 0x101), int(g / 0x101), int(b / 0x101)
}

func FindNearestColor(userRGB [3]byte) int {
	minDistance := math.Inf(1)
	var nearestColor int

	for _, hexColor := range Colors {
		predefinedRGB := [3]int{(hexColor >> 16) & 0xFF, (hexColor >> 8) & 0xFF, hexColor & 0xFF}
		distance := math.Sqrt(math.Pow(float64(int(userRGB[0])-predefinedRGB[0]), 2) +
			math.Pow(float64(int(userRGB[1])-predefinedRGB[1]), 2) +
			math.Pow(float64(int(userRGB[2])-predefinedRGB[2]), 2))

		if distance < minDistance {
			minDistance = distance
			nearestColor = hexColor
		}
	}

	return nearestColor
}

func MarkupText(str string) string {
	var sb strings.Builder

	state := "normal"
	var colorBuffer strings.Builder

	len := len(str)
	for i := 0; i < len; i++ {
		char := str[i]

		switch state {
		case "normal":
			if char == '[' {
				// Check if it's a escaped tag
				if i+1 < len && str[i+1] == '[' {
					Verbose("Escaped tag at index", fmt.Sprint(i+1))
					sb.WriteRune('[') // Add tag
					i++               // Go past escaped tag
					continue
				}

				state = "readingColor"
				colorBuffer.Reset() // Clear any previous color buffer
			} else {
				sb.WriteByte(char) // Write normal text
			}

		case "readingColor":
			if char == '[' {
				Verbose("Found unexpected character (ignoring)", string(char))
				continue
			}

			if str[i-1] == '[' && char == '/' {
				state = "readingReset"
			} else if char == ']' && !(i-1 >= 0 && str[i-1] == '[') {
				color := resolveColorCode(colorBuffer.String(), false)
				sb.WriteString(color) // Apply the resolved color code
				colorBuffer.Reset()
				state = "normal"
			} else {
				colorBuffer.WriteByte(char) // Build the color name
			}

		case "readingReset":
			if char == ']' && !(i-1 >= 0 && str[i-1] == '[') {
				sb.WriteString(RESET_COLOR) // Reset colors
			} else {
				Verbose("Found unexpected character (ignoring)", string(char))
				// Ignore

				//fmt.Printf("Encountered unexpected '%c' character when looking for closing tag ']' at index %d in snippet `%v'\n", char, i, getSurrounding(str, i, 5))
				//os.Exit(9)
			}
			state = "normal"
		}
	}

	// If we end in "readingColor" state, apply the last color code
	if state == "readingColor" {
		color := resolveColorCode(colorBuffer.String(), false)
		sb.WriteString(color)
	}

	finalOutput := sb.String()
	if argparse.Configuration.Escape {
		finalOutput = expandEscapeCodes(finalOutput)
	}

	return finalOutput
}

func resolveColorCode(colorStr string, resolvingBackground bool) string {

	// fmt.Printf("In:  `%v`\n", colorStr)
	if len(colorStr) == 0 || colorStr == "/" || argparse.Configuration.ForceColor == 1 {
		return RESET_COLOR
	}

	finalColor := ""
	if !resolvingBackground {
		segments := strings.Split(colorStr, " on ")
		if len(segments) == 2 {
			// Get the background color

			finalColor = resolveColorCode(segments[1], true)
			colorStr = segments[0]
		}
	}

	if argparse.Configuration.Verbose {
		Verbose("resolving tag", colorStr)
	}

	styleStr, colorStr := addExtraStyles(colorStr)

	if len(colorStr) != 0 && colorStr != "_" {
		// Get color
		if colorStr[0] == '#' && len(colorStr) == 7 {
			/*
			 *	Get a color by its hex value
			 */

			iColor, err := strconv.ParseInt(colorStr[1:], 16, 32)

			if err == nil {
				finalColor += Verbose("hex color", outputColorFromHex(int(iColor), resolvingBackground))
			}

		} else if strings.HasPrefix(colorStr, "rgb(") {
			/*
			 *	Get a color by its RGB formatted value
			 */

			// Parse RGB color (e.g. rgb(0, 0, 255))
			var r, g, b byte
			_, err := fmt.Sscanf(colorStr, "rgb(%d,%d,%d)", &r, &g, &b)

			if err == nil {
				finalColor += Verbose("rgb color", outputColorFromRgb([3]byte{r, g, b}, resolvingBackground))
			}

		} else {
			/*
			 *	Get a color by its name
			 */

			iColor, found := Colors[colorStr]

			if !found {
				Verbose("failed to find color", colorStr)
				return ""
			}

			finalColor += Verbose("named color", outputColorFromHex(iColor, resolvingBackground))
		}
	}

	if resolvingBackground {
		// Just give color codes, no escape char
		bOutput := fmt.Sprintf("%s;%s", styleStr, finalColor)
		bOutput = cleanSemicolons(bOutput)

		return Verbose("final background", bOutput)
	}

	// 'm' terminates the escape code, regardless of style data being present
	finalColor = fmt.Sprintf("%s;%s", styleStr, finalColor)
	finalColor = cleanSemicolons(finalColor)
	finalColor = fmt.Sprintf("\033[%sm", finalColor)
	return Verbose("final color", finalColor)
}

// Format a color from Hex to RGB string
func outputColorFromHex(iColor int, background bool) string {
	_r, _g, _b := IntHexToRGB(iColor)
	return outputColorFromRgb([3]byte{_r, _g, _b}, background)
}

// Format a color from RGB values ([3]byte) into a
func outputColorFromRgb(iColors [3]byte, background bool) string {
	colorType := 38
	if background {
		colorType = 48
	}

	r, g, b := IntHexToRGB(FindNearestColor(iColors))

	return fmt.Sprintf(";%d;2;%v;%v;%v", colorType, r, g, b)
}

// Replace keywords with styles other than colors
func addExtraStyles(str string) (string, string) {
	var sb strings.Builder
	color := ""

	// https://stackoverflow.com/questions/4842424/list-of-ansi-color-escape-sequences
	for _, item := range strings.Split(str, " ") {
		code, found := formatCodes[item]

		if found {
			sb.WriteString(code)
			Verbose("style item", item)
		} else {
			color = item
			Verbose("possible color item", item)
		}
	}

	outputStr := sb.String()
	if len(outputStr) != 0 {
		outputStr = outputStr[1:]
	}

	return Verbose("resolved styles", outputStr), Verbose("filtered color", color)
}

// Known (and supported) codes to manipulate terminal input
var formatCodes = map[string]string{
	"bold":       ";1",
	"faint":      ";2",
	"italic":     ";3",
	"underlined": ";4",
	"blinking":   ";5",
	"fblinking":  ";6", // "fast" blinking (some terminals render at the same speed as 'blinking')
	"swap":       ";7",
	"dunderline": ";21", // "double" underlined
	"overlined":  ";53",
}

func expandEscapeCodes(str string) string {
	for _, code := range escapeDictionary {
		str = strings.ReplaceAll(str, code, escapeDictionary[code])
	}

	return str
}

var escapeDictionary = map[string]string{
	"\\003": "\033",
	"\\x1b": "\x1b",
	"\\n":   "\n",
	"\\r":   "\r",
	"\\t":   "\t",
	"\\a":   "\a",
	"\\b":   "\b",
	"\\f":   "\f",
	"\\v":   "\v",
}

// Remove duplicate semicolons that are next to each other
func cleanSemicolons(text string) string {
	var result strings.Builder
	for i := 0; i < len(text); i++ {
		if text[i] == ';' {
			if i == 0 || text[i-1] != ';' {
				result.WriteByte(text[i])
			}
		} else {
			result.WriteByte(text[i])
		}
	}

	return strings.Trim(result.String(), ";")
}

// Write message to the console when --verbose is used
func Verbose(prefix string, in string) string {
	if argparse.Configuration.Verbose {
		caller := CallerName(1)
		fmt.Printf("[%s] %s: `%s`\n", caller, prefix, strings.ReplaceAll(in, "\033", "\\033"))
	}

	return in
}

func CallerName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return ""
	}

	f := runtime.FuncForPC(pc)
	if f == nil {
		return ""
	}

	name := f.Name()

	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}
