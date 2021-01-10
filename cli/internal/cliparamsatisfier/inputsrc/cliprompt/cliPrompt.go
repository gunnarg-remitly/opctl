package cliprompt

import (
	"fmt"

	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier/inputsrc"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/peterh/liner"
)

func New(
	cliOutput clioutput.CliOutput,
	inputs map[string]*model.Param,
) inputsrc.InputSrc {
	return cliPromptInputSrc{
		inputs:    inputs,
		cliOutput: cliOutput,
	}
}

// cliPromptInputSrc implements InputSrc interface by sourcing inputs from std in
type cliPromptInputSrc struct {
	inputs    map[string]*model.Param
	cliOutput clioutput.CliOutput
}

func (this cliPromptInputSrc) ReadString(
	inputName string,
) (*string, bool) {
	if param := this.inputs[inputName]; nil != param {
		var (
			isSecret    bool
			description string
			prompt      string
		)

		switch {
		case nil != param.Array:
			isSecret = param.Array.IsSecret
			description = param.Array.Description
			prompt = "array"
		case nil != param.Boolean:
			description = param.Boolean.Description
			prompt = "boolean"
		case nil != param.Dir:
			isSecret = param.Dir.IsSecret
			description = param.Dir.Description
			prompt = "directory"
		case nil != param.File:
			isSecret = param.File.IsSecret
			description = param.File.Description
			prompt = "file"
		case nil != param.Number:
			isSecret = param.Number.IsSecret
			description = param.Number.Description
			prompt = "number"
		case nil != param.Object:
			isSecret = param.Object.IsSecret
			description = param.Object.Description
			prompt = "object"
		case nil != param.Socket:
			isSecret = param.Socket.IsSecret
			description = param.Socket.Description
			prompt = "socket"
		case nil != param.String:
			isSecret = param.String.IsSecret
			description = param.String.Description
			prompt = "string"
		}
		prompt += ": "

		line := liner.NewLiner()
		defer line.Close()
		line.SetCtrlCAborts(true)

		if description != "" {
			this.cliOutput.Attention(
				fmt.Sprintf("input: \"%s\"\n%s", inputName, description),
			)
		} else {
			this.cliOutput.Attention(
				fmt.Sprintf("input: \"%s\"", inputName),
			)
		}

		// liner has inconsistent behavior if non empty prompt arg passed so use ""
		var (
			err    error
			rawArg string
		)
		if isSecret {
			rawArg, err = line.PasswordPrompt(prompt)
		} else {
			rawArg, err = line.Prompt(prompt)
		}
		if nil == err {
			return &rawArg, true
		}
	}

	return nil, false
}
