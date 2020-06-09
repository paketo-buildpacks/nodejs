# Nodejs implementation buildpacks architecture

## Summary

This document describes what implementation cnbs must exist, what their
function should be and how they should appear in the meta cnb order.

## Motivation

It looks like in the current state of the world, we have:

- a [node-engine](github.com/paketo-buildpacks/node-engine) cnb that provides
  `node` and `npm` on PATH.

- an [npm](github.com/paketo-buildpacks/npm) cnb that does all of these
  functions - install dependencies in the local `node_modules` folder, and sets
  a start command using `npm start`.

- a [yarn-install](github.com/paketo-buildpacks/yarn-install) cnb that does all
  of these functions - provide the yarn executable on PATH, install
  dependencies in the local `node_modules` folder, and sets a start command using
  `npm start`. The cnb is poorly named due to historical reasons

This above architecture does not go well with the implementation CNB philosophy
of "ask for your requirements at each stage, do one thing, and do it well".
This leads to additional dependencies ending up in the final app image
providing implicit rather than explicit behavior.

Further, this structure makes it difficult to address complex issues like
[#37](https://github.com/paketo-buildpacks/nodejs/issues/37)


## Proposal

Have 6 implementation CNBS:

- node-engine: provides `node` and `npm` (npm is the default package manager
  that comes together with node) executable on PATH
  - Requires none.

- npm-install: installs dependencies to `node_modules` and provides them.
  - Requires `node` during `build`

- yarn: provides `yarn` executable on PATH.
  - Requires none.

- yarn-install: installs dependencies to `node_modules` and provides them.
  - Requires `node` and `yarn` during `build`

- node-run: smart enough to come up with an appropriate start command to run
  either an npm app or yarn app or a vanilla node app. See Note [1].
  - Requires one of {node+yarn-install, node+npm-install, node} during `launch`


The above implementation buildpacks should be structured as follows in the nodejs language family buildpack.

```toml
[[order]]

  [[order.group]]
    id = "paketo-buildpacks/node-engine"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/yarn"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/yarn-install"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/node-run"
    version = ""

[[order]]

  [[order.group]]
    id = "paketo-buildpacks/node-engine"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/npm-install"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/node-run"
    version = ""


[[order]]

  [[order.group]]
    id = "paketo-buildpacks/node-engine"
    version = ""

  [[order.group]]
    id = "paketo-buildpacks/node-run"
    version = ""
```

## Notes

[1] This start command will not be either `npm start` or `yarn run/start` and
should construct a start command like "node [...]". This is because both `npm
start` and `yarn run` automatically shell out to node which creates problems
seen in [issue #37](https://github.com/paketo-buildpacks/nodejs/issues/37).
These commands do not seem very conducive for running production apps inside
docker containers.
