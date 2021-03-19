package docker

import (
	"bufio"
	"io"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
)

func NewStdOutWriteCloser(
	eventChannel chan model.Event,
	containerCall *model.ContainerCall,
	rootCallID string,
) io.WriteCloser {
	pr, pw := io.Pipe()
	go func() {
		// support lines up to 4MB
		reader := bufio.NewReaderSize(pr, 4e6)
		var err error
		var b []byte

		for {
			// chunk on newlines
			if b, err = reader.ReadBytes('\n'); len(b) > 0 {
				// always publish if len(bytes) read to ensure full stream sent; even under error conditions
				eventChannel <- model.Event{
					Timestamp: time.Now().UTC(),
					ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenTo{
						Data:        b,
						OpRef:       containerCall.OpPath,
						ContainerID: containerCall.ContainerID,
						RootCallID:  rootCallID,
					},
				}
			}

			if nil != err {
				return
			}
		}
	}()

	return pw
}
