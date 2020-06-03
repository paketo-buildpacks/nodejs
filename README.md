# Node.js Paketo Buildpack

## `gcr.io/paketo-buildpacks/nodejs`

The Node.js Paketo Buildpack provides a set of collaborating buildpacks that
enable the building of a Node.js-based application. These buildpacks include:
- [Node Engine CNB](https://github.com/paketo-buildpacks/node-engine)
- [Yarn Install CNB](https://github.com/paketo-buildpacks/yarn-install)
- [NPM CNB](https://github.com/paketo-buildpacks/npm)

The buildpack supports building simple Node applications or applications which
utilize either [NPM](https://www.npmjs.com/) or [Yarn](https://yarnpkg.com/)
for managing their dependencies. Support for each of these package managers is
mutually exclusive.
