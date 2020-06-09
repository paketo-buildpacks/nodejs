# Nodejs implementation buildpacks architecture

## Summary

This document describes what implementation cnbs must exist, what their
function should be and how they should appear in the meta cnb order.

## Motivation

It looks like in the current state of the world, we have:

- a [node-engine](github.com/paketo-buildpacks/node-engine) cnb that provides
  `node` on PATH.

- an [npm](github.com/paketo-buildpacks/npm) cnb that does all of these
  functions - provide the npm executable on PATH, install dependencies in the
  local `node_modules` folder, and sets a start command using `npm start`.

- a [yarn-install](github.com/paketo-buildpacks/yarn-install) cnb that does all
  of these functions - provide the yarn executable on PATH, install
  dependencies in the local `node_modules` folder, and sets a start command using
  `npm start`. The cnb is poorly named due to historical reasons

This above architecture does not go well with the implementation CNB philosophy
of "ask for your requirements, do one thing, and do it well".

Further, this structure makes it difficult to address complex issues like
[#37](https://github.com/paketo-buildpacks/nodejs/issues/37)


## Proposal

Have 6 implementation CNBS:

- node-engine: provides `node` executable on PATH

- npm: provides `npm` executable on PATH
- npm-install: installs dependencies to `node_modules` and provides them. Requires `node`.
- npm-run (another name?): comes up with a start command to run the app. Requires `node`. 

- yarn: provides `yarn` executable on PATH
- yarn-install: installs dependencies to `node_modules` and provides them. Requires `node`.
- yarn-run (another name?): comes up with a start command to run the app. Requires `node`. 


## Implementation

```toml
[[order]]

  [[order.group]]
    id = "paketo-buildpacks/"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/"
    version = ""

[[order]]

  [[order.group]]
    id = "paketo-buildpacks/"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/"
    version = ""

[[order]]

  [[order.group]]
    id = "paketo-buildpacks/"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/"
    version = ""
```
