# Gobis-buildpack over Cloud Foundry Experimental Multi-buildpack

This is a proof of concept to run a [gobis-server](https://github.com/orange-cloudfoundry/gobis-server) on top of an app. 

This is a fork from [multi-buildpack](https://github.com/cloudfoundry-incubator/multi-buildpack).

App will be start on other port than the one listening by cloud foundry and gobis-server will be run on this port and forward request to the app. This permit to enrich your app without pain through gobis.

## Why it's a fork ?

Multi-buildpack doesn't permit to another buildpack to hook other buildpack start commands but it does provide the logic how to do it. This only to show how we can do such a things to put an apache in front of an app, a zuul or of course a gobis-server.

## Usage

**Behaviour from multi-buildpack**:

- This buildpack looks for a `multi-buildpack.yml` file in the root of the application directory with structure:

```yaml
buildpacks:
  - https://github.com/cloudfoundry/go-buildpack
  - https://github.com/cloudfoundry/ruby-buildpack/releases/download/v1.6.23/ruby_buildpack-cached-v1.6.23.zip
  - https://github.com/cloudfoundry/nodejs-buildpack#v1.5.18
  - https://github.com/cloudfoundry/python-buildpack#develop
```

- The multi-buildpack will download + run all the buildpacks in this list in the specified order.

- It will use the app start command given by the last buildpack run, and gobis hook will wrap this command to start app on other port and run the gobis-server.

**Behaviour with gobis-server**:

- You can set a service to your app following [gobis-server on cloud foundry](https://github.com/orange-cloudfoundry/gobis-server#on-cloudfoundry)
- If you put a `gobis-config.yml` in your root app folder routes inside will be loaded too
- If you want to add [middleware_params](https://github.com/orange-cloudfoundry/gobis-middlewares) to the default route which will forward request to the wrapped app create a `gobis-params.yml` in your app root, example to add basic auth:

```yaml
basic_auth:
- user: myuser
  password: mypassword
```

## Details

- This will not work with system buildpacks. Ex. the following `multi-buildpack.yml` file will not work:

```yaml
buildpacks:
  - ruby_buildpack
  - go_buildpack
```

- The multi-buildpack will run the `bin/compile` and `bin/release` scripts for each specified buildpack.

## Disclaimer

It is not intended for production usage.
