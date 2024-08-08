package verbose

import (
	"fmt"
	"runtime"
	"strings"
)

var VerboseLoggingEnabled = false

// Write message to the console when --verbose is used
func Log(prefix string, in string) string {
	if VerboseLoggingEnabled {
		caller := CallerName(1)
		fmt.Printf("[%s] %s: `%s`\n",
			caller,
			prefix,
			strings.ReplaceAll(in, "\033", "\\033")) // Escape markup
	}

	return in
}

func CallerName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return ""
	}

	function := runtime.FuncForPC(pc)
	if function == nil {
		return ""
	}

	name := function.Name()

	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}
