package opstate

import (
	"testing"
)

func TestStripAnsi(t *testing.T) {
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

func TestStripAnsiToLength(t *testing.T) {
	ansiStr := "\033[1Atesting a string"
	stripped := stripAnsiToLength(ansiStr, 9)
	expected := "\033[1Atesting a"
	if stripped != expected {
		t.Errorf("stripped string isn't correct: expected `%s`, actual `%s`", expected, stripped)
	}
}

func TestStripAnsiToLength_escapeCodeInMid(t *testing.T) {
	ansiStr := "\033[1Atesting\033[0m a\033[1A string"
	stripped := stripAnsiToLength(ansiStr, 9)
	expected := "\033[1Atesting\033[0m a\033[1A"
	if stripped != expected {
		t.Errorf("stripped string isn't correct: expected `%s`, actual `%s`", expected, stripped)
	}
}

func TestStripAnsiToLength_noAnsi(t *testing.T) {
	ansiStr := "testing a string"
	stripped := stripAnsiToLength(ansiStr, 9)
	expected := "testing a"
	if stripped != expected {
		t.Errorf("stripped string isn't correct: expected `%s`, actual `%s`", expected, stripped)
	}
}
