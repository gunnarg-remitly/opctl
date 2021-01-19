package op

import (
	"testing"

	dataResolverFakes "github.com/opctl/opctl/cli/internal/dataresolver/fakes"
)

func TestOpNew(t *testing.T) {
	New(new(dataResolverFakes.FakeDataResolver))
}
