# :warning: Deprecated  in favor of [sidecars-buildpack](https://github.com/orange-cloudfoundry/sidecars-buildpack) :warning:


### Gobis buildpack

Gobis buildpack is a special buildpack to deploy a [gobis-server](https://github.com/orange-cloudfoundry/gobis-server) 
as a [sidecar](https://blog.davemdavis.net/2018/03/13/the-sidecar-pattern/).

It load every middleware you can found on https://github.com/orange-cloudfoundry/gobis-middlewares .

You simply need to have a latest cloud foundry which support multi-buildpack

This buildpack can't be used as a final buildpack and support stacks:
- cflinuxfs2
- cflinuxfs3
- windows2012R2
- windows2016
- windows

### Buildpack User Documentation


1. Add the buildpack as the first buildpack on your app manifest
2. Change the start command by `gobis-server --sidecar`
3. that's it to make pass all your traffic on gobis-server instead of your app (your app will listen on port `8081`, 
you can set internal domain on this port)

**Manifest example**:

```yaml
applications:
  - name: front
    buildpacks:
      - gobis_buildpack
      - staticfile_buildpack
    disk_quota: 1G
    command: gobis-server --sidecar # tips: you can use all cli params from gobis-server, add flag `--log-level debug` to enable debug mode for example
```

**Tips**: You can override start command for your app by creating a file named `Procfile` and add a `start` entry, e.g.:

```yaml
start: start-command-for-app
```

### Gobis configuration

Now your gobis-server is active you would like to add configuration on it.

create a file `route.yml` in an folder named `.gobis` (cmplete path: `.gobis/route.yml`)

**example of configuration**:

```yaml
# List of headers which should not be sent to upstream
sensitive_headers: []
# An url to an http proxy to make requests to upstream pass to this
http_proxy: ""
# An url to an https proxy to make requests to upstream pass to this
https_proxy: ""
# Force to never use proxy even proxy from environment variables
no_proxy: false
# By default response from upstream are buffered, it can be issue when sending big files
# Set to true to stream response
no_buffer: false
# Set to true to not send X-Forwarded-* headers to upstream
remove_proxy_headers: false
# Set to true to not check ssl certificates from upstream (not really recommended)
insecure_skip_verify: false
# Set to true to see errors on web page when there is a panic error on gobis
show_error: false
# Chain others routes in a routes
routes: ~
# Will forward directly to proxified route OPTIONS method without using middlewares
options_passthrough: false
middleware_params:
  cors:
    max_age: 12
    allowed_origins:
    - http://localhost
```

**Tips**: For `middleware_params` you can instead create a `.gobis/*-params.yml` file. 
All params inside all this files will be loaded as `middleware_params`, example for cors define in config example:

```yaml
# file is named `cors-params.yml`
cors:
  max_age: 12
  allowed_origins:
  - http://localhost
```

### How does it works ?

Gobis-server running as a sidecar will do the following:
1. Load `.gobis/route.yml` file and add only this route to gobis (`name` will be configured as `proxy-<app-name>` and `path` will be `/**`)
2. Look for all files named as follow `.gobis/*-params.yml` and load them as middleware params to be injected in route.
3. Get from env var `GOBIS_PORT` to know where app should listen
4. Create route url to `http://127.0.0.1:<previous found port>`
5. Look in `Procfile` if key `start` is found. Content is the custom command for real app that user want to override
6. Run default launcher from cloud foundry with start command given by user if exists 
to start real app with `PORT` env var override to previous found port 
7. Gobis-server will listening on port expected by cloud foundry and will reverse traffic to app beside 

### Building the Buildpack
To build this buildpack, run the following command from the buildpack's directory:

1. Source the .envrc file in the buildpack directory.
```bash
source .envrc
```
To simplify the process in the future, install [direnv](https://direnv.net/) which will automatically source .envrc when you change directories.

1. Install buildpack-packager
```bash
./scripts/install_tools.sh
```

1. Build the buildpack
```bash
buildpack-packager build
```

1. Use in Cloud Foundry
Upload the buildpack to your Cloud Foundry and optionally specify it by name

```bash
cf create-buildpack [BUILDPACK_NAME] [BUILDPACK_ZIP_FILE_PATH] 1
cf push my_app [-b BUILDPACK_NAME]
```

### Testing
Buildpacks use the [Cutlass](https://github.com/cloudfoundry/libbuildpack/cutlass) framework for running integration tests.

To test this buildpack, run the following command from the buildpack's directory:

1. Source the .envrc file in the buildpack directory.

```bash
source .envrc
```
To simplify the process in the future, install [direnv](https://direnv.net/) which will automatically source .envrc when you change directories.

1. Run unit tests

```bash
./scripts/unit.sh
```

1. Run integration tests

```bash
./scripts/integration.sh
```

More information can be found on Github [cutlass](https://github.com/cloudfoundry/libbuildpack/cutlass).

### Reporting Issues
Open an issue on this project

