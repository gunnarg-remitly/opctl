# Branch notes

## purpose

This branch (mini-opctl) is an exploration into fixing some of my main concerns
with opctl.

### reliability/complexity

I find opctl quite unreliable. The major issues are issues connecting to port
42224 and opctl not properly cleaning up after itself, resulting in docker
containers hanging around. The core architecture of opctl starts up a persistent
webserver on 42224 that actually manages running the core business logic of
opctl. The CLI then interacts with that webserver using http api calls and a
websocket connection to pipe output back to the user.

Since opctl's primary value to me is a reusable "op" runner, this webserver
feels very unnecessary. It doesn't maintain much internal state, and the state
it does maintain doesn't allow me to query and interact with running ops. The 
network interactions are what I suspect are the primary causes of reliability 
issues.

Because of this, I've entirely removed this webserver component of opctl in
favor of the CLI directly running core logic. This allows me to pass a
cancellable context through the entire call chain, and directly respond to
returned errors.

Opctl also uses a custom publisher/subscriber event bus internally, which
becomes pretty unnecessary once the centralized API is gone. I've replaced this
with a standard go channel that the events can be passed back through.

### line count

This project is _huge_ and can be difficult to work in. This ~removes checked-in
vendored code,~ the web UI for opctl, the JS sdk, and the react SDK.

The project also has many layers of abstraction, that I feel could be reduced
to make changes easier. I also think the code could be refactored to be more
idiomatic to the go language.

### usability

For complex ops, opctl makes it difficult to understand what's going on. I hope
to improve the output of the CLI tool to allow me to identify what produces
what output.

## features

This is a list of smaller features from this branch that I'd like to attempt to
incrementally migrate into the main branch, decoupled from the larger architecture
changes.

- Better CLI output (better = more readable, understandable, and transformable)
  - Label where output comes from (involves piping more context in events)
  - Clean up formatting (remove extraneous newlines, remove unnecessary formatting separators)
- Better error propagation and cleanup behavior
  - Ensure parallel call goroutines are waited on
  - Return and handle errors within the "Call" call stack
  - Ensure container cleanup won't happen with a cancelled context
- Move `ListDescendants` and `GetData` implementation to sdk core, instead of the api client
- Emit CallStarted events for skipped conditional branches [#859](https://github.com/opctl/opctl/pull/859)
- Remove custom pubsub?

## back to main

Long term, if the current "remote node" architecture is maintained, I'd like to
make it an opt-in feature to avoid needing a persistent process for day-to-day
local only use. This could be done by making the ApiClient and Core objects use
the same interface, which would also probably improve understandability of the
codebase, and would force better error propagation.

If the project focuses on a CLI runner model like this branch uses, we can still
support a persistent UI server by streaming events from the CLI to that server,
instead of the current model of round-tripping everything.

---

[![Build](https://github.com/opctl/opctl/workflows/Build/badge.svg?branch=main)](https://github.com/opctl/opctl/actions?query=workflow%3ABuild+branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/opctl/opctl)](https://goreportcard.com/report/github.com/opctl/opctl)
[![Coverage](https://codecov.io/gh/opctl/opctl/branch/master/graph/badge.svg)](https://codecov.io/gh/opctl/opctl)

> *Be advised: this project is currently at Major version zero. Per the
> semantic versioning spec: "Major version zero (0.y.z) is for initial
> development. Anything may change at any time. The public API should
> not be considered stable."*

# Documentation

see [website](https://opctl.io)

# Used By

These awesome companies use opctl. represent by adding yours (or one you know of) to the list!
- [Era](https://helloera.co)
- [Expedia](https://www.expedia.com)
- [Nintex](https://www.nintex.com)
- [ProKarma](https://prokarma.com/)
- [Remitly](https://www.remitly.com)
- [Samsung (SDS)](https://www.samsungsds.com)

# Support

join us on
[![Slack](https://img.shields.io/badge/slack-opctl-E01563.svg)](https://join.slack.com/t/opctl/shared_invite/zt-51zodvjn-Ul_UXfkhqYLWZPQTvNPp5w)
or [open an issue](https://github.com/opctl/opctl/issues)

# Releases

releases are versioned according to
[![semver 2.0.0](https://img.shields.io/badge/semver-2.0.0-brightgreen.svg)](http://semver.org/spec/v2.0.0.html)
and [tagged](https://git-scm.com/book/en/v2/Git-Basics-Tagging); see
[CHANGELOG.md](CHANGELOG.md) for release notes

# Contributing

see [CONTRIBUTING.md](CONTRIBUTING.md)


