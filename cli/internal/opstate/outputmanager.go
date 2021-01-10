package opstate

import (
	"fmt"
	"strings"

	"golang.org/x/term"
)

// OutputManager allows printing a "resettable" thing at the bottom of a stream
// of terminal output, when a tty is used
type OutputManager struct {
	lastHeight int
}

// NewOutputManager returns a new OutputManager
func NewOutputManager() OutputManager {
	return OutputManager{}
}

// Clear clears the last thing printed by this object
func (o *OutputManager) Clear() {
	// cursor to start of line (real big number)
	fmt.Print("\033[100000D")
	// clear line
	fmt.Print("\033[K")
	for i := 1; i < o.lastHeight; i++ {
		// move up one line
		fmt.Print("\033[1A")
		// clear line
		fmt.Print("\033[K")
	}
}

// Print prints the given string, with a preceding separator and width limited
// to the size of the terminal
func (o *OutputManager) Print(str string) error {
	w, _, err := term.GetSize(0)
	if err != nil {
		return err
	}
	lines := strings.Split(str, "\n")

	ruleWidth := 0
	for _, line := range lines {
		visualLen := countChars(stripAnsi(line))
		if visualLen > ruleWidth {
			ruleWidth = visualLen
		}
	}
	if ruleWidth > w {
		ruleWidth = w
	}

	fmt.Println(strings.Repeat("┄", ruleWidth))

	for i, line := range lines {
		withoutAnsi := stripAnsi(line)
		if countChars(withoutAnsi) > w {
			// - append an ellipsis to make it more obvious the line's being truncated
			// - remove _two_ chars, not just one for the ellipsis, because the cursor
			//   will take up an additional space and cause output issues
			// - append a "reset" code to prevent color wrapping to next line
			fmt.Print(stripAnsiToLength(line, w-2) + "…\033[0m")
		} else {
			fmt.Print(line)
		}
		if i < len(lines)-1 {
			fmt.Print("\n")
		}
	}

	o.lastHeight = len(lines) + 1
	return nil
}
