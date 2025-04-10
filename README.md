# Paketo Buildpack for Node.js

## `docker.io/paketobuildpacks/nodejs`

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

#### The Node.js buildpack is compatible with the following builder(s):

- [Paketo Jammy Full Builder](https://github.com/paketo-buildpacks/builder-jammy-full)
- [Paketo Jammy Base Builder](https://github.com/paketo-buildpacks/builder-jammy-base)

This buildpack also includes the following utility buildpacks:

- [Procfile CNB](https://github.com/paketo-buildpacks/procfile)
- [Environment Variables CNB](https://github.com/paketo-buildpacks/environment-variables)
- [Image Labels CNB](https://github.com/paketo-buildpacks/image-labels)
- [CA Certificates CNB](https://github.com/paketo-buildpacks/ca-certificates)
- [Node Run Script CNB](https://github.com/paketo-buildpacks/node-run-script)
- [Node Module Bill of Materials CNB](https://github.com/paketo-buildpacks/node-module-bom)

Check out the [Paketo Node.js docs](https://paketo.io/docs/buildpacks/language-family-buildpacks/nodejs/) for more information.

## Maintenance

### Go Module Versioning

Each buildpack is a Go module, and in the case of the Node.js buildpacks, we only maintain and support the latest versions, without providing support for any **previous** `major` or `minor` versions.  **Further the team does not currently commit to maintain the go major versions to in sync with the Semver versions used to publish releases. Keeping the go major versions up to date with the versions used to publish the buildpacks and consumed in the buildpack tomls is addhoc and based PRs being sumitted by the community**.
