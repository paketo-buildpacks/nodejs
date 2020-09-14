# Remove Tini Buildpack

## Summary

The inclusion of `tini` in this buildpack is no longer required as the
`npm-start` and `yarn-start` buildpacks have moved away from its use. We should
remove references to the `tini` buildpack.

## Motivation

The integration tests for `npm-start` and `yarn-start` have shown that using
`tini` to manage processes in concert with `npm` and `yarn` is unreliable.
Those buildpacks have been factored in a [Yarn Start
RFC0002](https://github.com/paketo-buildpacks/yarn-start/blob/main/rfcs/0002-reimplement-start-command.md)
and [NPM Start RFC
0002](https://github.com/paketo-buildpacks/npm-start/blob/main/rfcs/0002-reimplement-start-command.md)
to remove `tini` from their launch process command.

The language-family buildpack should be updated to reflect these changes.

## Proposal

Remove the `tini` buildpack from the Yarn and NPM buildpack group orderings.
When complete, the `buildpack.toml` should have the following structure:

```toml
[[order]]

  [[order.group]]
    id = "paketo-buildpacks/node-engine"

  [[order.group]]
    id = "paketo-buildpacks/yarn"

  [[order.group]]
    id = "paketo-buildpacks/yarn-install"

  [[order.group]]
    id = "paketo-buildpacks/yarn-start"

[[order]]

  [[order.group]]
    id = "paketo-buildpacks/node-engine"

  [[order.group]]
    id = "paketo-buildpacks/npm-install"

  [[order.group]]
    id = "paketo-buildpacks/npm-start"

[[order]]

  [[order.group]]
    id = "paketo-buildpacks/node-engine"

  [[order.group]]
    id = "paketo-buildpacks/node-start"
```
