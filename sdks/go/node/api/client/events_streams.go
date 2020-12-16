package client

import (
	"context"
	"fmt"

	"github.com/opctl/opctl/sdks/go/model"
)

func (c client) GetEventStream(
	ctx context.Context,
	req *model.GetEventStreamReq,
) (<-chan model.Event, error) {
	eventStream, errChan := c.core.GetEventStream(ctx, req)

	go func() {
		for {
			err := <-errChan
			if err != nil {
				fmt.Printf("error received: %v\n", err)
			}
		}
	}()

	return eventStream, nil
}
