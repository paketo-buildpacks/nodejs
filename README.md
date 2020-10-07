# Node.js Paketo Buildpack

## `gcr.io/paketo-buildpacks/nodejs`

The Node.js Paketo Buildpack provides a set of collaborating buildpacks that
enable the building of a Node.js-based application. These buildpacks include:
- [Node Engine CNB](https://github.com/paketo-buildpacks/node-engine)
- [Yarn CNB](https://github.com/paketo-buildpacks/yarn)
- [Yarn Install CNB](https://github.com/paketo-buildpacks/yarn-install)
- [NPM Install CNB](https://github.com/paketo-buildpacks/npm-install)
- [Yarn Start CNB](https://github.com/paketo-buildpacks/yarn-start)
- [NPM Start CNB](https://github.com/paketo-buildpacks/npm-start)
- [Node Start CNB](https://github.com/paketo-buildpacks/node-start)

The buildpack supports building/running simple Node applications or applications
which utilize either [NPM](https://www.npmjs.com/) or [Yarn](https://yarnpkg.com/)
for managing their dependencies. Support for each of these package managers is
mutually exclusive.

Usage examples can be found in the
[`samples` repository under the `nodejs` directory](https://github.com/paketo-buildpacks/samples/tree/main/nodejs).
