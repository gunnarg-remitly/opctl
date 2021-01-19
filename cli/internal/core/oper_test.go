package core

import (
	"testing"

	dataResolverFakes "github.com/opctl/opctl/cli/internal/dataresolver/fakes"
)

func TestNewOper(t *testing.T) {
	if newOper(new(dataResolverFakes.FakeDataResolver)).Op() == nil {
		t.Error("Oper should provide Op")
	}
}
