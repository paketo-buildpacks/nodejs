# Nodejs implementation buildpacks architecture

## Summary

This document describes what implementation CNB must exist, what their
function should be and how they should appear in the meta CNB order.

## Motivation

It looks like in the current state of the world, we have:

- a [node-engine](github.com/paketo-buildpacks/node-engine) CNB that provides
  `node` and `npm` on PATH.

- an [npm](github.com/paketo-buildpacks/npm) CNB that does all of these
  functions - install dependencies in the local `node_modules` folder, and sets
  a start command using `npm start`.

- a [yarn-install](github.com/paketo-buildpacks/yarn-install) CNB that does all
  of these functions - provide the yarn executable on PATH, install
  dependencies in the local `node_modules` folder, and sets a start command using
  `npm start`. The CNB is poorly named due to historical reasons

This above architecture does not go well with the implementation CNB philosophy
of "ask for your requirements at each stage, do one thing, and do it well".
This leads to additional dependencies ending up in the final app image
providing implicit rather than explicit behavior.

Further, this structure makes it difficult to address complex issues like
[#37](https://github.com/paketo-buildpacks/nodejs/issues/37)


## Proposal

Have 8 implementation CNBs:

- node-engine:
  - Provides `node` or {`node`, `npm`} (npm is the default package manager
  that comes together with node) executable on PATH
  - Requires none.

- npm-install: installs dependencies to `node_modules` and provides them.
  - Provides `node_modules`
  - Requires {`node`, `npm`} during `build`

- yarn: provides `yarn` executable on PATH.
  - Provides `yarn`
  - Requires none.

- yarn-install: installs dependencies to `node_modules` and provides them.
  - Provides `node_modules`
  - Requires {`node`, `yarn`} during `build`

- tini: provides the [tini](https://github.com/krallin/tini) process manager,
  which will be used in our start commands. This is needed to allow signals to
  be passed to the node process spawned by the following start command CNBs
  npm-start and yarn-start
  - Provides `tini`
  - Requires none.
  See [issue #37](https://github.com/paketo-buildpacks/nodejs/issues/37) for motivation.


- npm-start: sets up a start command that uses `tini` and `npm`
  - Provides none
  - Requires {`node`, `npm`, `node_modules`, `tini`} during `launch`
  - Example start command would be `tini -g -- npm start`

- yarn-start: sets up a start command that uses `tini` and `yarn`
  - Provides none
  - Requires {`node`, `yarn`, `node_modules`, `tini`} during `launch`

- node-start: sets up a start command that uses `node`
  - Provides none
  - Requires {`node`} during `launch`
  - Example start command would be `node server.js`

The above implementation buildpacks should be structured as follows in the nodejs language family buildpack.

```toml
[[order]]

  [[order.group]]
    id = "paketo-buildpacks/node-engine"

  [[order.group]]
    id = "paketo-buildpacks/yarn"

  [[order.group]]
    id = "paketo-buildpacks/yarn-install"

  [[order.group]]
    id = "paketo-buildpacks/tini"

  [[order.group]]
    id = "paketo-buildpacks/yarn-start"

[[order]]

  [[order.group]]
    id = "paketo-buildpacks/node-engine"

  [[order.group]]
    id = "paketo-buildpacks/npm-install"

  [[order.group]]
    id = "paketo-buildpacks/tini"

  [[order.group]]
    id = "paketo-buildpacks/npm-start"

[[order]]

  [[order.group]]
    id = "paketo-buildpacks/node-engine"

  [[order.group]]
    id = "paketo-buildpacks/node-start"
```

**Revision History**

* (09/09/2020) Edit: Fix npm-install's api - it was conceptualized to be the same as
  that of yarn-install with npm in place of yarn.
