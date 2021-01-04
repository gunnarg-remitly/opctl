package opfile

import (
	"bytes"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/opctl/opctl/sdks/go/model"
)

// Unmarshal validates and unmarshals an "op.yml" file
func Unmarshal(
	opRef string,
	opFileBytes []byte,
) (*model.OpSpec, error) {
	// 1) ensure valid
	errs := Validate(opFileBytes)
	if len(errs) > 0 {
		messageBuffer := bytes.NewBufferString("opspec syntax error:\n")
		messageBuffer.WriteString(opRef)
		for _, validationError := range errs {
			messageBuffer.WriteString(fmt.Sprintf("\n- %v", validationError.Error()))
		}
		return nil, fmt.Errorf("%v", messageBuffer.String())
	}

	// 2) build
	opFile := model.OpSpec{}
	return &opFile, yaml.Unmarshal(opFileBytes, &opFile)
}
