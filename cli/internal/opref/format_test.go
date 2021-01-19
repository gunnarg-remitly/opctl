package opref

import (
	"crypto/rand"
	"testing"

	. "github.com/onsi/gomega"
)

func TestSimpleFormatter(t *testing.T) {
	g := NewGomegaWithT(t)

	objectUnderTest := SimpleOpFormatter{}

	rand.Read()
	g.Expect(objectUnderTest.FormatOpRef(""))
}

func TestOpFormatter() {

}
