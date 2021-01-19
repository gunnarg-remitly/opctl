package clicolorer

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"github.com/fatih/color"
)

//counterfeiter:generate -o fakes/colorer.go . CliColorer

type colorFunc func(string) string

// CliColorer allows adding color to text before displaying through the CLI.
// It has cross platform terminal support.
type CliColorer interface {
	// DisableColor disables any coloring of CLI output, for use in environments
	// that won't support ANSI like color escape codes.
	DisableColor()

	Muted(s string) string
	Attention(string) string
	Error(string) string
	Info(string) string
	Success(s string) string
}

// New creates a new CliColorer
func New() CliColorer {
	mutedCliColorer := color.New(color.Faint)
	errorCliColorer := color.New(color.FgHiRed, color.Bold)
	attentionCliColorer := color.New(color.FgHiYellow, color.Bold)
	infoCliColorer := color.New(color.FgHiCyan, color.Bold)
	successCliColorer := color.New(color.FgHiGreen, color.Bold)

	return &cliColorer{
		mutedCliColorer:     mutedCliColorer,
		errorCliColorer:     errorCliColorer,
		attentionCliColorer: attentionCliColorer,
		infoCliColorer:      infoCliColorer,
		successCliColorer:   successCliColorer,
	}
}

type cliColorer struct {
	mutedCliColorer     *color.Color
	errorCliColorer     *color.Color
	attentionCliColorer *color.Color
	infoCliColorer      *color.Color
	successCliColorer   *color.Color
}

func (c *cliColorer) DisableColor() {
	c.mutedCliColorer.DisableColor()
	c.errorCliColorer.DisableColor()
	c.attentionCliColorer.DisableColor()
	c.infoCliColorer.DisableColor()
	c.successCliColorer.DisableColor()
}

func (c cliColorer) Muted(s string) string {
	return c.mutedCliColorer.SprintFunc()(s)
}

func (c cliColorer) Error(s string) string {
	return c.errorCliColorer.SprintFunc()(s)
}

func (c cliColorer) Attention(s string) string {
	return c.attentionCliColorer.SprintfFunc()(s)
}

func (c cliColorer) Info(s string) string {
	return c.infoCliColorer.SprintFunc()(s)
}

func (c cliColorer) Success(s string) string {
	return c.successCliColorer.SprintFunc()(s)
}
