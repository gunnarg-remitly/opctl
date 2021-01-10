package opstate

import (
	"testing"
)

func TestAnsiEscape(t *testing.T) {
	str := "test"
	ansiStr := "\033[1Atest"
	withoutAnsi := stripAnsi(ansiStr)

	if str != withoutAnsi {
		t.Error("stripped string is not equal to original")
	}

	str = "◉ ⠴ ./test"
	ansiStr = "◉ ⠴ [1m./test[0m"
	withoutAnsi = stripAnsi(ansiStr)

	if str != withoutAnsi {
		t.Error("stripped string is not equal to original")
	}
}
