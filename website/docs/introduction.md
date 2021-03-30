---
title: Introduction
sidebar_label: Introduction
---

Opctl is a high level system for running and plumbing data through containerized tasks.

An op (operation) is defined with an opspec (operation specification), a yaml file-based schema.

## Use cases

- Avoid implementation-specific dependency management
- Share tasks across different platforms (e.g. CI and local development)
- Reuse common components between different systems (e.g. fetching configuration for different projects)

## CLI

The opctl CLI runs your ops locally. See [reference docs](reference/cli.md) for full details.

![Example opctl CLI output](/img/cli-output.png)

## Opspec

Opspec is a language designed to portably and fully define ops. See [reference docs](reference/opspec/index.md) for full details.

It features:
- Containers as first class citizens
- Serial and parallel looping and execution
- Conditional execution
- Variables and scoping
- Array, boolean, dir, file, number, object, socket, and string data types
- Explicit inputs/outputs with type specific constraints
- Composition and re-use of ops
- Versioning via git tags
- Type coercion
- Declarative dependencies between calls
