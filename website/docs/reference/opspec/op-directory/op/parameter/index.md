---
sidebar_label: Index
title: Parameter [object]
---

An object defining a parameter of an operation; i.e. a value that is passed into or out of it's scope.

## Properties
- must have exactly one of
  - [array](array.md)
  - [boolean](boolean.md)
  - [dir](dir.md)
  - [file](file.md)
  - [number](number.md)
  - [object](object.md)
  - [socket](socket.md)
  - [string](string.md)

## Example
```yaml
name: example
description: an example op
inputs:
  example-input:
    string:
      default: "a default value"
run:
  container:
    image: { ref: 'alpine' }
    cmd: ['echo', $(example-input)]
```
