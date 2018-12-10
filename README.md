[![Build Status](https://build.vio.sh/buildStatus/icon?job=vapor-ware/synse-snmp-plugin/master)](https://build.vio.sh/blue/organizations/jenkins/vapor-ware%2Fsynse-snmp-plugin/activity)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-snmp-plugin.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-snmp-plugin?ref=badge_shield)
![GitHub release](https://img.shields.io/github/release/vapor-ware/synse-snmp-plugin.svg)

# Synse SNMP Plugin
A general-purpose SNMP plugin for [Synse Server][synse-server].

## Plugin Support
### Outputs
Outputs should be referenced by name. A single device can have more than one instance
of an output type. In the table below, a value of `-` indicates that there is no value
set for that field.

| Name | Description | Unit | Precision | Scaling Factor |
| ---- | ----------- | ---- | --------- | -------------- |
| `current` | An output type for current readings | amps | 3 | - |
| `frequency` | An output type for frequency readings | hertz | 3 | - |
| `identity` | An output type for identity readings | - | - | - |
| `va.power` | An output type for power readings, in volt-amperes | volt-ampere | 3 | - |
| `watts.power` | An output type for power readings, in watts | watts | 3 | - |
| `status` | An output type for status readings | - | - | - |
| `temperature` | An output type for temperature readings | degrees celsius | 3 | - |
| `voltage` | An output type for voltage readings | volts | 3 | - |

### Device Handlers
Device Handlers should be referenced by name.

| Name | Description | Read | Write | Bulk Read |
| ---- | ----------- | ---- | ----- | --------- |
| `current` | A handler for the SNMP OIDs that report current | ✓ | ✗ | ✗ |
| `frequency` | A handler for the SNMP OIDs that report frequency | ✓ | ✗ | ✗ |
| `identity` | A handler for the SNMP-identity device | ✓ | ✗ | ✗ |
| `power` | A handler for the SNMP OIDs that report power | ✓ | ✗ | ✗ |
| `status` | A handler for the SNMP-status device | ✓ | ✗ | ✗ |
| `temperature` | A handler for the SNMP OIDs that report temperature | ✓ | ✗ | ✗ |
| `voltage` | A handler for the SNMP OIDs that report voltage | ✓ | ✗ | ✗ |

### Write Values
The SNMP plugin does not currently support writing to any devices.


## Getting Started
### Getting the Plugin
You can get the Synse SNMP plugin by:

0. Cloning this repo, setting up the project dependencies, and building the binary or docker image
   ```bash
   # setup the project
   $ make setup

   # build the binary
   $ make build

   # build the docker image
   $ make docker
   ```
0. Pulling a pre-built docker image from [DockerHub][plugin-dockerhub]
   ```bash
   $ docker pull vaporio/snmp-plugin
   ```
0. Downloading a pre-built binary from the latest [release][plugin-release].

### Running the Plugin
If you are using the Docker image (recommended):
```bash
$ docker run vaporio/snmp-plugin
```

If you are using the binary directly:
```bash
# The name of the plugin binary may differ depending on how it is saved.
$ ./plugin
```

In either case, the plugin should run but you should not see any devices configured
and you should see errors in the logs saying that various configurations were not found.
The next section describes how to configure the plugin.

## Configuring the Plugin
Plugin and device configurations are described in detail in the [SDK Configuration Documentation][sdk-config-docs].

For your own deployment, you will need to provide your own plugin config, `config.yml`.
The SNMP plugin dynamically generates the device configuration records from a MIB walk,
so there is no need to specify device instance configuration files here.

For a simple example of how dynamic registration works, see the SDK's
[dynamic registration example][dynamic-reg-example].

### plugin config
After reading through the docs linked above for the plugin config, take a look at the
[example plugin config](config.yml). This can be used as a reference. To specify your
own SNMP devices, you will need to list them under the `dynamicRegistration` field, e.g.
```yaml
dynamicRegistration:
  config:
    - model: PXGMS UPS + EATON 93PM
      version: v3
      endpoint: 127.0.0.1
      port: 1024
      userName: simulator
      authenticationProtocol: SHA
      authenticationPassphrase: auctoritas
      privacyProtocol: AES
      privacyPassphrase: privatus
      contextName: public
```

Once you have your plugin config, you can either mount it into the container at `/plugin/config.yml`,
or mount it anywhere, e.g. `/tmp/cfg/config.yml`, and specify that path in the plugin config
override environment variable, `PLUGIN_CONFIG=/tmp/cfg`.

## Supported MIBs

* [UPS-MIB][ups-mib-rfc]


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

## Feedback
Feedback for this plugin, or any component of the Synse ecosystem, is greatly appreciated!
If you experience any issues, find the documentation unclear, have requests for features,
or just have questions about it, we'd love to know. Feel free to open an issue for any
feedback you may have.

## Contributing
We welcome contributions to the project. The project maintainers actively manage the issues
and pull requests. If you choose to contribute, we ask that you either comment on an existing
issue or open a new one.

## License
The Synse SNMP Plugin, and all other components of the Synse ecosystem, is released under the
[GPL-3.0](LICENSE) license.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-snmp-plugin.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-snmp-plugin?ref=badge_large)


[synse-server]: https://github.com/vapor-ware/synse-server
[plugin-dockerhub]: https://hub.docker.com/r/vaporio/snmp-plugin
[plugin-release]: https://github.com/vapor-ware/synse-snmp-plugin/releases
[sdk-config-docs]: http://synse-sdk.readthedocs.io/en/latest/user/configuration.html
[dynamic-reg-example]: https://github.com/vapor-ware/synse-sdk/tree/master/examples/dynamic_registration
[ups-mib-rfc]: https://tools.ietf.org/html/rfc1628
[issues]: https://github.com/vapor-ware/synse-snmp-plugin/issues