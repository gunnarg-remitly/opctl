package opstate

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestDotLoadingSpinner(t *testing.T) {
	// arrange
	var objectUnderTest DotLoadingSpinner

	// act
	l := objectUnderTest.String()

	// assert
	_, size := utf8.DecodeRuneInString(l)
	assert.Equal(t, 3, size)
	assert.Equal(t, size, len(l))
}
