package core

import (
	"context"
	"os"
	"testing"

	"github.com/opctl/opctl/cli/internal/clioutput"
	cliOutputFakes "github.com/opctl/opctl/cli/internal/clioutput/fakes"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
)

func TestNewCore(t *testing.T) {
	// arrange
	objectUnderTest, err := New(
		context.Background(),
		new(cliOutputFakes.FakeCliOutput),
		clioutput.SimpleOpFormatter{},
		local.NodeCreateOpts{
			DataDir: os.TempDir(),
		},
	)

	// assert
	if err != nil {
		t.Error(err)
	}
	if objectUnderTest.Auth() == nil {
		t.Error("core should provide Auth")
	}
	if objectUnderTest.Op() == nil {
		t.Error("core should provide Op")
	}
}
