[![Build Status](https://build.vio.sh/buildStatus/icon?job=vapor-ware/synse-snmp-plugin/master)](https://build.vio.sh/blue/organizations/jenkins/vapor-ware%2Fsynse-snmp-plugin/activity)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-snmp-plugin.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-snmp-plugin?ref=badge_shield)
![GitHub release](https://img.shields.io/github/release/vapor-ware/synse-snmp-plugin.svg)

# Synse SNMP Plugin

A general-purpose SNMP plugin for [Synse Server][synse-server].

## Plugin Support

### Outputs

Outputs are referenced by name.

Below is a table detailing the Outputs defined by the plugin. For an accounting of Outputs
built-in to the SDK, see [builtins.go](https://github.com/vapor-ware/synse-sdk/blob/master/sdk/output/builtins.go).

| Name | Type | Description | Unit | Precision |
| ---- | ---- | ----------- | ---- | --------- |
| `volt-ampere` | power | An output type for power readings, in volt-amperes | volt-ampere (VA) | 3 |
| `identity` | identity | An output type for identity readings | - | - |

### Device Handlers

Device Handlers should be referenced by name.

| Name | Description | Read | Write | Bulk Read |
| ---- | ----------- | ---- | ----- | --------- |
| `current` | A handler for the SNMP OIDs which report current | ✓ | ✗ | ✗ |
| `frequency` | A handler for the SNMP OIDs which report frequency | ✓ | ✗ | ✗ |
| `identity` | A handler for the SNMP-identity device(s) | ✓ | ✗ | ✗ |
| `power` | A handler for the SNMP OIDs which report power | ✓ | ✗ | ✗ |
| `status` | A handler for the SNMP-status device(s) | ✓ | ✗ | ✗ |
| `temperature` | A handler for the SNMP OIDs which report temperature | ✓ | ✗ | ✗ |
| `voltage` | A handler for the SNMP OIDs which report voltage | ✓ | ✗ | ✗ |

> **Note**: The SNMP plugin does not currently support writing to any device. This support
> may be added in the future.

## Getting Started

### Getting the Plugin

It is recommended to run the plugin in a Docker container. You can pull the image from
[DockerHub][plugin-dockerhub]:

```
docker pull vaporio/snmp-plugin
```

You can also download a plugin binary from the latest [release][plugin-release].

### Running the Plugin

If you are using the Docker image (recommended):

```bash
$ docker run vaporio/snmp-plugin
```

If you are using the plugin binary:

```bash
# The name of the plugin binary may differ depending on how it is built/downloaded.
$ ./plugin
```

In either case, the plugin should run but you should not see any devices configured
and you should see errors in the logs saying that various configurations were not found.
The next section describes how to configure the plugin.

## Configuring the Plugin

Device and plugin configuration are described in the [Synse SDK Documentation][sdk-docs].

For your deployment, you will need to provide your own plugin config, `config.yml`.
The SNMP plugin dynamically generates the device configuration records from a MIB walk,
so there is no need to specify device instance configuration files here.

For a simple example of how dynamic registration works, see the SDK's
[dynamic registration example][dynamic-reg-example].

### Plugin Config

After reading through the docs linked above for the plugin config, take a look at the
[example plugin config](example/config.yml). This may be used as a reference. To specify your
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

## Supported MIBs

* [UPS-MIB][ups-mib-rfc]

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
[sdk-docs]: http://synse-sdk.readthedocs.io/en/latest/user/configuration.html
[dynamic-reg-example]: https://github.com/vapor-ware/synse-sdk/tree/master/examples/dynamic_registration
[ups-mib-rfc]: https://tools.ietf.org/html/rfc1628
