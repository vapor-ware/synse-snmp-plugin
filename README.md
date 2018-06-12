# Vapor SNMP Plugin (Internal)

The internal SNMP plugin ports over the SNMP functionality from the OpenDCRE
implementation.

## Architecture
The SNMP plugin uses the [Synse Go SDK][go-sdk]. It follows a similar pattern to other plugins
where the plugin-specific read/write methods (see [plugin.go][plugin-main]) dispatch to
device-specific handlers based on the model of the requested device. This allows each device to
have its own arbitrary implementation.

### Supported MIBs

* [UPS-MIB][ups-mib-rfc]

Plugins are used in conjunction with Synse Server; they provide the backend data which
Synse Server makes available to any upstream API user.

See the `synse-plugins-internal/deploy/` directory for examples on how to run Synse
Server with plugins. Examples exist for both docker-compose based deployments
and kubernetes based deployments. Furthermore, working examples using the SNMP plugin
can be found there for both unix socket networking configurations and TCP networking
configurations.

### Configuration
There are three type of configuration for plugins:
- plugin configuration
- device prototype configuration
- device instance configuration

The device prototype configuration is found in `config/proto/` and defines base attributes
of a device. For each prototype configuration in `config/proto/`, there should be a matching
device implementation in `devices/`. The device prototype configurations are built into the
plugin and define the set of devices that the plugin supports.

The device instance configuration is found in `config/device/` and defines specific devices
that the plugin should communicate with. The configurations here are plugin-dependent and for
the most part specify how to reach that device (e.g. which serial port to use, etc). The repo
contains sample device configurations, but these are not built into the plugin. They should be
specified at runtime.

The plugin configuration is found in `plugin.yml`. It provides options for the behavior of the
plugin, such as whether it should run in debug mode, the timeout interval, read/write queue
sizes, etc.


## Deployment
Generally, there are three ways to deploy a plugin:
- directly on the host
- built-in to the Synse Server image
- as a standalone image

The primary intended usage is as a standalone image. Below describes how to deploy
the plugin via Docker as well as to a Kubernetes cluster.

### Docker
**TODO**

### Kubernetes
**TODO**


## Developing
The Makefile provides targets to simplify the development workflow. For all targets,
use `make help`

### Building
There are two builds that are supported - building the Go binary locally, and building
the docker image.

To build the Go binary
```
make build
```

To build the Docker image
```
make docker
```

Note that when building the docker image, it will first build the Go binary for linux.

### Linting
Prior to committing code, the source should be linted. Linting is built into the CI
pipeline, so if linting fails, the code will not be merged in.
```
make lint
```

### Testing
Prior to committing code, the source should be tested. Tests should be added for any new
component.
```
make test
```


## Troubleshooting
### Debugging
The plugin can be run in debug mode for additional logging. This can be done by setting
```yaml
debug: true
```
in the plugin configuration yaml (`plugin.yml`)

### Bugs / Issues
If you find a bug or experience an issue, open a [new issue][issues] and provide as much context
information as possible. What happened? What was expected to happen? How do they differ?
Are there logs that document what happened? etc.



[plugin-main]: https://github.com/vapor-ware/synse-plugins-internal/blob/master/i2c/plugin.go
[go-sdk]: https://github.com/vapor-ware/synse-sdk
[issues]: https://github.com/vapor-ware/synse-snmp-plugin/issues
[ups-mib-rfc]: https://tools.ietf.org/html/rfc1628
