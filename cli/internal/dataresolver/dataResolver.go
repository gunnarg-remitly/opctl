package dataresolver

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/data/fs"
	"github.com/opctl/opctl/sdks/go/data/node"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/core"
)

// DataResolver resolves packages
//counterfeiter:generate -o fakes/dataResolver.go . DataResolver
type DataResolver interface {
	Resolve(
		ctx context.Context,
		dataRef string,
		pullCreds *model.Creds,
	) (model.DataHandle, error)
}

func New(
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	core core.Core,
) DataResolver {
	return _dataResolver{
		cliParamSatisfier: cliParamSatisfier,
		core:              core,
		os:                ios.New(),
	}
}

type _dataResolver struct {
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier
	core              core.Core
	os                ios.IOS
}

func (dtr _dataResolver) Resolve(
	ctx context.Context,
	dataRef string,
	pullCreds *model.Creds,
) (model.DataHandle, error) {
	cwd, err := dtr.os.Getwd()
	if nil != err {
		return nil, err
	}

	fsProvider := fs.New(
		filepath.Join(cwd, ".opspec"),
		cwd,
	)

	for {
		opDirHandle, err := data.Resolve(
			ctx,
			dataRef,
			fsProvider,
			node.New(dtr.core, pullCreds),
		)

		var isAuthError bool
		switch err.(type) {
		case model.ErrDataProviderAuthorization:
			isAuthError = true
		case model.ErrDataProviderAuthentication:
			isAuthError = true
		}

		switch {
		case nil == err:
			return opDirHandle, err
		case isAuthError:
			// auth errors can be fixed by supplying correct creds so don't give up; prompt
			argMap, err := dtr.cliParamSatisfier.Satisfy(
				cliparamsatisfier.NewInputSourcer(
					dtr.cliParamSatisfier.NewCliPromptInputSrc(credsPromptInputs),
				),
				credsPromptInputs,
			)
			if nil != err {
				return nil, err
			}

			// save providedArgs & re-attempt
			pullCreds = &model.Creds{
				Username: *(argMap[usernameInputName].String),
				Password: *(argMap[passwordInputName].String),
			}
			continue
		default:
			// uncorrectable error.. give up
			return nil, fmt.Errorf("Unable to resolve pkg '%v'; error was %v", dataRef, err.Error())
		}

	}

}

const (
	usernameInputName = "username"
	passwordInputName = "password"
)

var (
	credsPromptInputs = map[string]*model.Param{
		usernameInputName: {
			String: &model.StringParam{
				Description: "username used to auth w/ the pkg source",
				Constraints: map[string]interface{}{
					"MinLength": 1,
				},
			},
		},
		passwordInputName: {
			String: &model.StringParam{
				Description: "password used to auth w/ the pkg source",
				Constraints: map[string]interface{}{
					"MinLength": 1,
				},
				IsSecret: true,
			},
		},
	}
)
