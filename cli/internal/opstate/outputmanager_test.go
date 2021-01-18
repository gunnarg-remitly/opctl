package opstate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOutputManager(t *testing.T) {
	objectUnderTest := NewOutputManager()

	assert.NotNil(t, objectUnderTest)
	w, err := objectUnderTest.getWidth()
	assert.Error(t, err, "in tests, terminal width isn't available")
	assert.Equal(t, -1, w)

	assert.Error(t, objectUnderTest.Print(""))
	objectUnderTest.Clear()
}

func TestOutputManagerShortLines(t *testing.T) {
	// arrange
	var buff bytes.Buffer
	objectUnderTest := OutputManager{
		getWidth: func() (int, error) { return 80, nil },
		out:      &buff,
	}

	// act
	err := objectUnderTest.Print(`testing
the drinks should just be like a glass of water
inspired by space jam and who knows`)

	// assert
	assert.Nil(t, err)
	expected := `┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄
testing
the drinks should just be like a glass of water
inspired by space jam and who knows`
	assert.Equal(t, expected, buff.String())
	assert.Equal(t, 4, objectUnderTest.lastHeight)
}

func TestOutputManagerLongLines(t *testing.T) {
	// arrange
	var buff bytes.Buffer
	objectUnderTest := OutputManager{
		getWidth: func() (int, error) { return 80, nil },
		out:      &buff,
	}
	longLine := strings.Repeat("-", 80)
	longerLine := strings.Repeat("-", 81)

	// act
	err := objectUnderTest.Print(fmt.Sprintf(`testing
%s
%s
testing3`, longLine, longerLine))

	// assert
	assert.Nil(t, err)
	expected := fmt.Sprintf(`┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄
testing
--------------------------------------------------------------------------------
------------------------------------------------------------------------------…%s
testing3`, "\033[0m")
	assert.Equal(t, expected, buff.String())
	assert.Equal(t, 5, objectUnderTest.lastHeight)
}

func TestOutputManagerClearing(t *testing.T) {
	// arrange
	var buff bytes.Buffer
	objectUnderTest := OutputManager{
		getWidth: func() (int, error) { return 80, nil },
		out:      &buff,
	}

	// act
	objectUnderTest.Print(`testing
the drinks should just be like a glass of water
inspired by space jam and who knows`)
	ioutil.ReadAll(&buff)
	objectUnderTest.Clear()

	// assert
	expected := "\x1b[80D\x1b[K\x1b[1A\x1b[K\x1b[1A\x1b[K\x1b[1A\x1b[K"
	assert.Equal(t, expected, buff.String())
	assert.Equal(t, 4, objectUnderTest.lastHeight)
}
