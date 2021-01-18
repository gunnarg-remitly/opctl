package opstate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripAnsi(t *testing.T) {
	str := "test"
	ansiStr := "\033[1Atest"
	withoutAnsi := stripAnsi(ansiStr)

	assert.Equal(t, str, withoutAnsi, "stripped string is not equal to original")
}

func TestStripAnsi_noAnsi(t *testing.T) {
	str := "◉ ⠴ ./test"
	ansiStr := "◉ ⠴ [1m./test[0m"
	withoutAnsi := stripAnsi(ansiStr)

	assert.Equal(t, str, withoutAnsi, "stripped string is not equal to original")
}

func TestStripAnsiToLength(t *testing.T) {
	ansiStr := "\033[1Atesting a string"
	stripped := stripAnsiToLength(ansiStr, 9)
	expected := "\033[1Atesting a"

	assert.Equal(t, expected, stripped)
}

func TestStripAnsiToLength_escapeCodeInMid(t *testing.T) {
	ansiStr := "\033[1Atesting\033[0m a\033[1A string"
	stripped := stripAnsiToLength(ansiStr, 9)
	expected := "\033[1Atesting\033[0m a\033[1A"

	assert.Equal(t, expected, stripped)
}

func TestStripAnsiToLength_noAnsi(t *testing.T) {
	ansiStr := "testing a string"
	stripped := stripAnsiToLength(ansiStr, 9)
	expected := "testing a"

	assert.Equal(t, expected, stripped)
}
