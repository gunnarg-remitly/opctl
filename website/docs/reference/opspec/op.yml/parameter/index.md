---
sidebar_label: Overview
title: Parameters
---

A parameter describes a value passed into or out of an op through it's inputs or outputs.

A parameter is an object with a single key declaring its type.

- [array](array.md)
- [boolean](boolean.md)
- [dir](dir.md)
- [file](file.md)
- [number](number.md)
- [object](object.md)
- [socket](socket.md)
- [string](string.md)

The value of each type of key is a further object. All types have a few properties in common:

## Common properties

### `default`

A literal value of the type of the parameter used as the default value of the variable created by the parameter.

### `description`

_required_

A human friendly description of the parameter, written as a [markdown string](markdown.md).
