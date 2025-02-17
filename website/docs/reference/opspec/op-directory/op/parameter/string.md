---
title: String Parameter [object]
---

An object defining a parameter which accepts a [string typed value](../../../types/string.md).

## Properties:
- may have:
  - [constraints](#constraints)
  - [default](#default)
  - [isSecret](#issecret)

### constraints
A [JSON Schema v4 [object]](https://tools.ietf.org/html/draft-wright-json-schema-00) defining constraints on the parameters value.

#### default
A string to use as the value of the parameter when no argument is provided.

#### isSecret
A boolean indicating if the value of the parameter is secret. This will cause it to be hidden in UI's for example.