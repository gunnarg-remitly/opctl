package opstate

import (
	"testing"
	"unicode/utf8"
)

func TestDotLoadingSpinner(t *testing.T) {
	var loader DotLoadingSpinner
	l := loader.String()
	if _, size := utf8.DecodeRuneInString(l); size != 3 {
		t.Error("loading spinner didn't behave as expected")
	}
	if len(l) != 3 {
		t.Error("loading spinner didn't behave as expected")
	}
}
