package colors

import (
	"gecko/argparse"

	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
)

const (
	NO_COLOR = "\033[0m"
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

func FindNearestColor(userRGB [3]byte, colorMap map[string]int) int {
	minDistance := math.Inf(1)
	var nearestColor int

	for _, hexColor := range colorMap {
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

func MarkupText(pcolors map[string]int, conf argparse.Config) string {
	var sb strings.Builder
	str := conf.TextInput

	state := "normal"
	var colorBuffer strings.Builder

	len := len(str)
	for i := 0; i < len; i++ {
		char := str[i]

		switch state {
		case "normal":
			if char == '[' && !(i-1 >= 0 && str[i-1] == '[') {
				state = "readingColor"
				colorBuffer.Reset() // Clear any previous color buffer
				//i++                 // Skip the opening '['
			} else {
				sb.WriteByte(char) // Write normal text
			}

		case "readingColor":
			if str[i-1] == '[' && char == '/' {
				state = "readingReset"
				//i += 2 // Skip "[/"
			} else if char == ']' {
				color := resolveColorCode(colorBuffer.String(), pcolors, conf)
				sb.WriteString(color) // Apply the resolved color code
				//sb.WriteByte(str[i+1])
				colorBuffer.Reset()
				state = "normal"
			} else {
				colorBuffer.WriteByte(char) // Build the color name
			}

		case "readingReset":
			if char == ']' {
				sb.WriteString(NO_COLOR) // Reset colors
				//state = "normal"
			} else {
				// Ignore
				//fmt.Printf("Encountered unexpected '%c' character when looking for closing tag ']' at index %d in snippet `%v'\n", char, i, getSurrounding(str, i, 5))
				//os.Exit(9)
			}
			state = "normal"
		}
	}

	// If we end in "readingColor" state, apply the last color code
	if state == "readingColor" {
		color := resolveColorCode(colorBuffer.String(), pcolors, conf)
		sb.WriteString(color)
	}

	return sb.String()
}

func getSurrounding(line string, from int, count int) string {
	textLength := len(line)

	if from < 0 {
		from = 0
	} else if from >= textLength {
		from = textLength - 1
	}

	count = min(count, textLength)
	start := max(0, from-count)
	end := min(from+count+1, textLength)

	return line[start:end]
}

func resolveColorCode(colorStr string, pcolors map[string]int, conf argparse.Config) string {

	// fmt.Printf("In:  `%v`\n", colorStr)
	if len(colorStr) == 0 || colorStr == "/" || conf.ForceColor == 1 {
		return NO_COLOR
	}

	if colorStr[0] == '#' && len(colorStr) == 7 {
		// Parse hex color (e.g., #ff04aa)
		value, err := strconv.ParseInt(colorStr[1:], 16, 32)

		if err == nil {
			return ret(fmt.Sprintf("\033[38;5;%dm", value))
		}

	} else if strings.HasPrefix(colorStr, "rgb(") {
		// Parse RGB color (e.g., rgb(0, 0, 255))
		var r, g, b int
		_, err := fmt.Sscanf(colorStr, "rgb(%d, %d, %d)", &r, &g, &b)

		if err == nil {
			return ret(fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b))
		}

	} else {
		// Delegate to your color name mapping logic
		sColor, found := pcolors[colorStr]

		if !found {
			return colorStr // just return the input if not a color
		}

		_r, _g, _b := IntHexToRGB(sColor)

		// fmt.Printf("cols: 'r:%d', 'g:%d', 'b:%d'\n", _r, _g, _b)
		r, g, b := IntHexToRGB(FindNearestColor([3]byte{_r, _g, _b}, pcolors))

		return ret(fmt.Sprintf("\033[38;2;%v;%v;%vm", r, g, b))
	}

	return NO_COLOR // Default to a color reset
}

func ret(in string) string {
	// fmt.Printf("out: `%s`\n", in)
	return in
}
