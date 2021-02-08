# Support running npm scripts

## Proposal

Add support for users to provide [npm scripts](https://docs.npmjs.com/cli/v6/configuring-npm/package-json#scripts)
to be executed after installing modules. Users will provide them via a
build-time environment variable and the execution will handled by a
`node-run-script` buildpack.

## Motivation

There are cases where users like the buildpack to run npm scripts in addition
to installing dependencies. Many frontend web application frameworks like
Angular, React, Vue.js etc. use node package managers like npm and yarn to
generate static files like html, js and css. The application workflow using
these frameworks involves generating a `package.json` which typically contains
a script under `build` (or another custom event name) that is run to generate
these static assets.

For e.g., a user building a simple Angular app would use the `ng new` command
to [generate a package.json](https://angular.io/guide/npm-packages#packagejson)
that looks like follows:
```
 "scripts": {
    "ng": "ng",
    "start": "ng serve",
    "build": "ng build",
    "test": "ng test",
    "lint": "ng lint",
    "e2e": "ng e2e"
  },
```

Here the command to build static resources are in the `build` script event. The
user then runs the command `npm run-script build` (or its alias `npm run build`)
to invoke this script which generates the static files in the
`<app>/dist` directory. This is a pattern we see in other frameworks as well.
Here are a few app examples:
* [React app](https://github.com/facebook/create-react-app/tree/v4.0.1)
* [Vue.js app](https://github.com/gothinkster/vue-realworld-example-app)
* [Augur.js app](https://github.com/AugurProject/augur-app/tree/v1.16.11)


As of now, the only way for the user to have these executed with buildpacks is
by editing the `package.json` and copying the command from the custom event to
a `preinstall` or `postinstall` script event that [npm and yarn recognizes](https://docs.npmjs.com/cli/v6/using-npm/scripts#npm-install)
as part of the `npm install`.  Many users do not wish to make changes to their
app source code, and there has been requests to support execution of these
commands by the buildpacks
([e.g.](https://github.com/paketo-buildpacks/yarn/issues/59)).

## Implementation

This RFC proposes a new `node-run-script` buildpack that runs the lifecycle
events provided by the user via build-time environment variable
`$BP_NODE_RUN_SCRIPTS`. The value of the variable should be a comma separated
list of events listed in the app's `package.json`.

e.g.
```shell
$BP_NODE_RUN_SCRIPTS="build,custom-script"
```

From the contents of the app directory and the buildpack plan entries, the
buildpack should pick the correct package manager (npm, yarn) to execute the
script.

The buildpack will figure in the Node.js language family buildpack order as
follows:

```
[[order]]
  [[order.group]]
    id = "paketo-buildpacks/node-engine"

  [[order.group]]
    id = "paketo-buildpacks/yarn"

  [[order.group]]
    id = "paketo-buildpacks/yarn-install"

  [[order.group]]
    id = "paketo-buildpacks/node-run-script"
    optional = true

  [[order.group]]
    id = "paketo-buildpacks/yarn-start"
  ...

[[order]]
  [[order.group]]
    id = "paketo-buildpacks/node-engine"

  [[order.group]]
    id = "paketo-buildpacks/npm-install"

  [[order.group]]
    id = "paketo-buildpacks/node-run-script"
    optional = true

  [[order.group]]
    id = "paketo-buildpacks/npm-start"
```

### Open questions

1. ~~Web framework examples mostly demonstrate cases where a `postinstall` event
   would be useful. What are some real world cases (if any) where `preinstall`
   would be used?~~ Just support running scripts after installing modules.

1. ~~Do NPM and Yarn need separate env variables or should they read the same set
   of env vars (e.g. `$BP_NODE_{PRE,POST}INSTALL`)?~~ Consolidate.

1. ~~Should the implementation be moved away from the `{npm,yarn}-install`
   buildpacks to separate npm run-script buildpack(s)?~~ Yes.

## Logistics

* A new buildpack repository named `node-run-script` should be created in the
  `paketo-buildpacks` organization under the Node.js subteam.
