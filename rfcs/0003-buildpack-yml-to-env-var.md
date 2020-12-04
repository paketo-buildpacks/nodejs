# Buildpack.yml to Environment Variables

## Proposal

Migrate to using environment variables to do all buildpack configuration and
get rid of `buildpack.yml`.

## Motivation

There are several reasons for making this switch.
1. There is already an existing RFC that proposes moving away from
   `buildpack.yml` as a configuration tool.
1. Environment variables appears to be the standard for configuration in other
   buildpack ecosystems such as Google Buildpacks and Heroku as well as the
   Paketo Java buildpacks. Making this change will align the buildpack with the
   rest of the buildpack ecosystem.
1. There is native support to pass environment variables to the buildpack
   either on a per run basis or by configuration that can be checked into
   source control, in the form of `project.toml`.
1. Go and Dotnet-core buildpacks are moving to this strategy.

## Implementation
The proposed environment variables are as follows:

### Existing configurations

#### BP_NODE_VERSION
```shell
$BP_NODE_VERSION="~10"
```
This will replace the following structure in `buildpack.yml`:
```yaml
nodejs:
  version: ~10
```

#### BP_NODE_OPTIMIZE_MEMORY
```shell
$BP_NODE_OPTIMIZE_MEMORY=true
```
This will replace the following structure in `buildpack.yml`:
```yaml
nodejs:
  optimize-memory: true
```

### Proposed additions

#### BP_NODE_PROJECT_PATH
```shell
$BP_NODE_PROJECT_PATH=./src/nodejs
```
If set, this relative path will be considered as the root of the app to be
built by nodejs buildpacks.

There is no current `buildpack.yml` support for this.

#### BP_LAUNCHPOINT
```shell
$BP_LAUNCHPOINT=./src/nodejs/customserver.js
```
Only applicable when using the `node-start` buildpack. This sets the file to
call as argument to node.
There is no current `buildpack.yml` support for this.

## Deprecation Strategy
In order to facilitate a smooth transition from `buildpack.yml`, the buildpack
should support both configuration options with environment variables taking
priority or `buildpack.yml` until the 1.0 release of the buildpack. The
buildpack will detect whether or not the application has a `buildpack.yml` and
print a warning message which will include links to documentation on how to
upgrade and how to run builds with environment variable configuration. After
1.0, having a `buildpack.yml` will cause a detection failure and with a link to
the same documentation. This behavior will only last until the next minor
release of the buildpack after which point there will no longer be and error
but `buildpack.yml` will not be supported.


## Related RFCs

* [Go build](https://github.com/paketo-buildpacks/go-build/pull/76)
* [Dotnet core](https://github.com/paketo-buildpacks/dotnet-core/pull/364)

## Source Material
* [Google buildpack configuration](https://github.com/GoogleCloudPlatform/buildpacks#language-idiomatic-configuration-options)
* [Paketo Java configuration](https://paketo.io/docs/buildpacks/language-family-buildpacks/java)
* [Heroku configuration](https://github.com/heroku/java-buildpack#customizing)
