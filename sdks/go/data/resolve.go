package data

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/pkg/errors"
)

// Resolve "dataRef" from "providers" in order
//
// expected errs:
//  - ErrDataProviderAuthentication on authentication failure
//  - ErrDataProviderAuthorization on authorization failure
//  - ErrDataRefResolution on resolution failure
func Resolve(
	ctx context.Context,
	dataRef string,
	providers ...model.DataProvider,
) (
	model.DataHandle,
	error,
) {
	var errs []error
	for _, src := range providers {
		handle, err := src.TryResolve(ctx, dataRef)
		if nil != err {
			errs = append(errs, errors.Wrap(err, src.Label()))
		} else if nil != handle {
			return handle, nil
		}
	}

	messageBuffer := bytes.NewBufferString(fmt.Sprintf("unable to resolve op \"%s\":", dataRef))
	for _, err := range errs {
		errStr := err.Error()
		parts := strings.Split(errStr, "\n")
		if len(parts) > 1 {
			for i, part := range parts {
				prefix := " "
				if i == 0 {
					prefix = "-"
				}
				messageBuffer.WriteString(fmt.Sprintf("\n%s %s", prefix, part))
			}
		} else {
			messageBuffer.WriteString(fmt.Sprintf("\n- %v", err))
		}
	}
	return nil, fmt.Errorf(messageBuffer.String())
}
