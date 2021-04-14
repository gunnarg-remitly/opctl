package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/opctl/opctl/sdks/go/model"
)

//counterfeiter:generate -o internal/fakes/imagePuller.go . imagePuller
type imagePuller interface {
	Pull(
		ctx context.Context,
		containerCall *model.ContainerCall,
		imagePullCreds *model.Creds,
		imageRef string,
		rootCallID string,
		eventChannel chan model.Event,
	) error
}

func newImagePuller(
	dockerClient dockerClientPkg.CommonAPIClient,
) imagePuller {
	return _imagePuller{
		dockerClient,
	}
}

type _imagePuller struct {
	dockerClient dockerClientPkg.CommonAPIClient
}

func (ip _imagePuller) Pull(
	ctx context.Context,
	containerCall *model.ContainerCall,
	imagePullCreds *model.Creds,
	imageRef string,
	rootCallID string,
	eventChannel chan model.Event,
) error {
	_, _, err := ip.dockerClient.ImageInspectWithRaw(ctx, imageRef)
	if err == nil {
		eventChannel <- model.Event{
			Timestamp: time.Now().UTC(),
			ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenTo{
				Data:        []byte(fmt.Sprintf("Skipping image pull: %s\n", imageRef)),
				OpRef:       containerCall.OpPath,
				ContainerID: containerCall.ContainerID,
				RootCallID:  rootCallID,
			},
		}

		return nil
	}

	imagePullOptions := types.ImagePullOptions{}
	if imagePullCreds != nil &&
		imagePullCreds.Username != "" &&
		imagePullCreds.Password != "" {
		var err error
		imagePullOptions.RegistryAuth, err = constructRegistryAuth(
			imagePullCreds.Username,
			imagePullCreds.Password,
		)
		if err != nil {
			return err
		}
	}

	imagePullResp, err := ip.dockerClient.ImagePull(
		ctx,
		imageRef,
		imagePullOptions,
	)
	if err != nil {
		return err
	}
	defer imagePullResp.Close()

	stdOutWriter := NewStdOutWriteCloser(eventChannel, containerCall, rootCallID)
	defer stdOutWriter.Close()

	dec := json.NewDecoder(imagePullResp)
	for {
		var jm jsonmessage.JSONMessage
		if err = dec.Decode(&jm); err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}
		jm.Display(stdOutWriter, false)
	}
}
