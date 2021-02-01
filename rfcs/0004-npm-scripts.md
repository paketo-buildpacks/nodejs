# Support running npm scripts

## Proposal

Add support for users to provide [npm scripts](https://docs.npmjs.com/cli/v6/configuring-npm/package-json#scripts)
to be executed at specific life cycle events, like before and after installing modules.
Users will provide them via build-time environment variables as defined in this RFC.

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
The proposed environment variables are as follows:

#### BP_NPM_PREINSTALL
```shell
$BP_NPM_PREINSTALL="build,custom-script"
```
A comma separated list of life cycle events listed in the app's `package.json`
that the `npm-install` buildpack must execute **before** it installs modules
using `npm`.

#### BP_NPM_POSTINSTALL
```shell
$BP_NPM_POSTINSTALL="build,custom-script"
```
A comma separated list of life cycle events listed in the app's `package.json`
that the `npm-install` buildpack must execute **after** it installs modules using `npm`.

#### BP_YARN_PREINSTALL
```shell
$BP_YARN_PREINSTALL="build,custom-script"
```
A comma separated list of life cycle events listed in the app's `package.json`
that the `yarn-install` buildpack must execute **before** it installs modules
using `yarn`.

#### BP_YARN_POSTINSTALL
```shell
$BP_YARN_POSTINSTALL="build,custom-script"
```
A comma separated list of life cycle events listed in the app's `package.json`
that the `yarn-install` buildpack must execute **after** it installs modules
using `yarn`.

### Order of execution

The life cycle events specified via environment variables is in addition to,
and must not interfere with the life cycle operation order of `npm install` or
other commands.

In the following case:

```
# Build-time environment variables:
$BP_NPM_PREINSTALL="env-preinstall-script"
$BP_NPM_POSTINSTALL="env-postinstall-script"

# package.json:
{
  "scripts" : {
    "preinstall" : "packagejson-preinstall-script",
    "install" : "packagejson-install-script",
    "postinstall" : "packagejson-postinstall-script"
    ...
    ...
  }
}
```

The following must be the order of execution of events by the npm-install buildpack:
```
env-preinstall-script
packagejson-preinstall-script
packagejson-install-script
packagejson-postinstall-script
env-postinstall-script
```

The *env-&ast;-script*s should be explicitly executed by the buildpacks and the
*packagejson-&ast;-script*s are implicitly executed as part of the buidpack's
[npm-install process](https://docs.npmjs.com/cli/v6/using-npm/scripts#npm-install).

### Open questions

1. Web framework examples mostly demonstrate cases where a `postinstall` event
   would be useful. What are some real world cases (if any) where `preinstall`
   would be used?

1. Do NPM and Yarn need separate env variables or should they read the same set
   of env vars (e.g. `$BP_NODE_{PRE,POST}INSTALL`)?

1. Should the implementation be moved away from the `{npm,yarn}-install`
   buildpacks to separate npm run-script buildpack(s)?
