package node

import (
	"errors"

	"github.com/opctl/opctl/cli/internal/nodeprovider"
)

// DisabledNoder is a non-interactive noder replacement
type DisabledNoder struct{}

// Node returns a non-interactive node replacement
func (DisabledNoder) Node() Node {
	return disabledNode{}
}

type disabledNode struct{}

func (disabledNode) Create(opts nodeprovider.NodeOpts) error {
	return errors.New("not initialized with a creatable node")
}

func (disabledNode) Kill() error {
	return errors.New("not initialized with a killable node")
}
