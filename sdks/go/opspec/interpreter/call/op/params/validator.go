package params

import (
	"bytes"
	"fmt"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/param"
)

// Validate validates values for/against params
func Validate(
	values map[string]*model.Value,
	params map[string]*model.Param,
) error {

	paramErrMap := map[string][]error{}
	for paramName, paramValue := range params {
		errs := param.Validate(values[paramName], paramValue)
		if len(errs) > 0 {
			paramErrMap[paramName] = errs
		}
	}

	if len(paramErrMap) > 0 {
		// return error w/ fancy formatted msg
		messageBuffer := bytes.NewBufferString("validation error(s):")
		for paramName, errs := range paramErrMap {
			for _, err := range errs {
				messageBuffer.WriteString(fmt.Sprintf("\n- %v: %v", paramName, err.Error()))
			}
		}
		return fmt.Errorf(messageBuffer.String())
	}

	return nil
}
