package main

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/appdataspec/sdk-golang/appdatapath"
	mow "github.com/jawher/mow.cli"
	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/cli/internal/clioutput"
	corePkg "github.com/opctl/opctl/cli/internal/core"
	"github.com/opctl/opctl/cli/internal/model"
	op "github.com/opctl/opctl/sdks/go/opspec"
)

//counterfeiter:generate -o internal/fakes/cli.go . cli
type cli interface {
	Run(args []string) error
}

// newCorer allows swapping out corePkg.New for unit tests
type newCorer func(
	ctx context.Context,
	cliOutput clioutput.CliOutput,
	containerRuntime,
	datadirPath string,
) (corePkg.Core, error)

func newCli(
	ctx context.Context,
	newCorer newCorer,
) (cli, error) {
	cli := mow.App(
		"opctl",
		"Opctl is a free and open source distributed operation control system.",
	)
	cli.Version("v version", version)

	perUserAppDataPath, err := appdatapath.New().PerUser()
	if nil != err {
		return nil, err
	}

	dataDir := cli.String(
		mow.StringOpt{
			Desc:   "Path of dir used to store opctl data",
			EnvVar: "OPCTL_DATA_DIR",
			Name:   "data-dir",
			Value:  filepath.Join(perUserAppDataPath, "miniopctl"),
		},
	)

	cliOutput, err := clioutput.New(clicolorer.New(), *dataDir, os.Stderr, os.Stdout)
	if err != nil {
		return nil, err
	}

	containerRuntime := cli.String(
		mow.StringOpt{
			Desc:   "Runtime for opctl containers",
			EnvVar: "OPCTL_CONTAINER_RUNTIME",
			Name:   "container-runtime",
			Value:  "docker",
		},
	)

	core, err := newCorer(ctx, cliOutput, *containerRuntime, *dataDir)
	if err != nil {
		return nil, err
	}

	exitWith := func(successMessage string, err error) {
		if err == nil {
			if successMessage != "" {
				cliOutput.Success(successMessage)
			}
			mow.Exit(0)
		} else {
			cliOutput.Error(err.Error())
			if re, ok := err.(*corePkg.RunError); ok {
				mow.Exit(re.ExitCode)
			} else {
				mow.Exit(1)
			}
		}
	}

	noColor := cli.BoolOpt("nc no-color", false, "Disable output coloring")

	cli.Before = func() {
		if *noColor {
			cliOutput.DisableColor()
		}
	}

	cli.Command("auth", "Manage auth for OCI image registries", func(authCmd *mow.Cmd) {
		authCmd.Command(
			"add", "Add auth for an OCI image registry",
			func(addCmd *mow.Cmd) {
				addCmd.Spec = "RESOURCES [ -u=<username> ] [ -p=<password> ]"

				resources := addCmd.StringArg("RESOURCES", "", "Resources this auth applies to in the form of a host or host/path.")
				username := addCmd.StringOpt("u username", "", "Username")
				password := addCmd.StringOpt("p password", "", "Password")

				addCmd.Action = func() {
					exitWith("", core.Auth().Add(ctx, *resources, *username, *password))
				}
			})
	})

	cli.Command("ls", "List operations (only valid ops will be listed)", func(lsCmd *mow.Cmd) {
		const dirRefArgName = "DIR_REF"
		lsCmd.Spec = fmt.Sprintf("[%v]", dirRefArgName)
		dirRef := lsCmd.StringArg(dirRefArgName, op.DotOpspecDirName, "Reference to dir ops will be listed from")

		lsCmd.Action = func() {
			core.Ls(ctx, *dirRef)
		}
	})

	cli.Command("op", "Manage ops", func(opCmd *mow.Cmd) {
		opCmd.Command("create", "Create an op", func(createCmd *mow.Cmd) {
			path := createCmd.StringOpt("path", op.DotOpspecDirName, "Path the op will be created at")
			description := createCmd.StringOpt("d description", "", "Op description")
			name := createCmd.StringArg("NAME", "", "Op name")

			createCmd.Action = func() {
				core.Op().Create(*path, *description, *name)
			}
		})

		opCmd.Command("install", "Install an op", func(installCmd *mow.Cmd) {
			path := installCmd.StringOpt("path", op.DotOpspecDirName, "Path the op will be installed at")
			opRef := installCmd.StringArg("OP_REF", "", "Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)")
			username := installCmd.StringOpt("u username", "", "Username used to auth w/ the pkg source")
			password := installCmd.StringOpt("p password", "", "Password used to auth w/ the pkg source")

			installCmd.Action = func() {
				exitWith("", core.Op().Install(ctx, *path, *opRef, *username, *password))
			}
		})

		opCmd.Command("validate", "Validate an op", func(validateCmd *mow.Cmd) {
			opRef := validateCmd.StringArg("OP_REF", "", "Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)")

			validateCmd.Action = func() {
				core.Op().Validate(ctx, *opRef)
			}
		})
	})

	cli.Command("run", "Start and wait on an op", func(runCmd *mow.Cmd) {
		args := runCmd.StringsOpt("a", []string{}, "Explicitly pass args to op in format `-a NAME1=VALUE1 -a NAME2=VALUE2`")
		argFile := runCmd.StringOpt("arg-file", filepath.Join(op.DotOpspecDirName, "args.yml"), "Read in a file of args in yml format")
		opRef := runCmd.StringArg("OP_REF", "", "Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)")

		runCmd.Action = func() {
			exitWith("", core.Run(ctx, *opRef, &model.RunOpts{Args: *args, ArgFile: *argFile}))
		}
	})

	cli.Command("self-update", "Update opctl", func(selfUpdateCmd *mow.Cmd) {
		channel := selfUpdateCmd.StringOpt("c channel", "stable", "Release channel to update from (either `stable`, `alpha`, or `beta`)")
		selfUpdateCmd.Action = func() {
			exitWith(core.SelfUpdate(*channel))
		}
	})

	return cli, nil
}
