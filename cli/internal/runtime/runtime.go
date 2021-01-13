package runtime

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/node/core/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime/docker"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime/k8s"
)

// New parses a string argument to create a container runtime
func New(containerRuntime string) (containerruntime.ContainerRuntime, error) {
	switch containerRuntime {
	case "k8s":
		return k8s.New()
	case "docker":
		return docker.New()
	}
	return nil, fmt.Errorf("unsupported runtime: %s", containerRuntime)
}
