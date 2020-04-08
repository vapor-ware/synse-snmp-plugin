[![Build Status](https://build.vio.sh/buildStatus/icon?job=vapor-ware/synse-snmp-plugin/master)](https://build.vio.sh/blue/organizations/jenkins/vapor-ware%2Fsynse-snmp-plugin/activity)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-snmp-plugin.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-snmp-plugin?ref=badge_shield)
![GitHub release](https://img.shields.io/github/release/vapor-ware/synse-snmp-plugin.svg)

# Synse SNMP Plugin

An SNMP plugin for [Synse Server][synse-server].

> Note: While there are plans to make this plugin more general-purposed in
> the future, it is currently limited in scope for what it supports. See the
> [Supported MIBs](#supported-mibs) section below.

## Getting Started

### Getting

You can install the SNMP plugin via a [release](https://github.com/vapor-ware/synse-snmp-plugin/releases)
binary or via Docker image

```
docker pull vaporio/snmp-plugin
```

If you wish to use a development build, fork and clone the repo and build the plugin
from source.

### Running

The SNMP plugin requires the SNMP-enabled servers it will communicate with to be configured.
As such, running the plugin without additional configuration will cause it to fail. As an
example of how to configure and get started with running the SNMP plugin, a simple example
deployment exists which runs runs Synse Server, the SNMP plugin, and a basic SNMP simulator.

To run it:

```bash
docker-compose up -d
```

You can then use Synse's HTTP API or the [Synse CLI][synse-cli] to query Synse for plugin data.

## SNMP Plugin Configuration

Plugin and device configuration are described in detail in the [Synse SDK Documentation][sdk-docs].

When deploying, you will need to provide your own plugin configuration (`config.yaml`)
with dynamic configuration defined. This is how the SNMP plugin knows about which servers
to communicate with. Dynamic configuration generates the devices at plugin startup from
the results of a MIB walk, so there is no need to specify device instance configuration
for the plugin.

As an example:

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

### Dynamic Registration Options

Below are the fields that are expected in each of the dynamic registration items.
If no default is specified (`-`), the field is required.

| Field                    | Description | Default |
| ------------------------ | ----------- | ------- |
| model                    | The model of the UPS. (Currently only supports models starting with "PXGMS UPS") | `-` |
| version                  | The SNMP protocol version. (Currently only "v3" is supported) | `-` |
| endpoint                 | The endpoint of the SNMP server to connect to. | `-` |
| port                     | The UDP port to connect to. | `-` |
| userName                 | The SNMP username. | `-` |
| authenticationProtocol   | The SNMP authentication protocol. (Supported: MD5, SHA) | `-` |
| authenticationPassphrase | The passphrase for authentication. | `-` |
| privacyProtocol          | The SNMP privacy protocol. (Supported: AES, DES) | `-` |
| privacyPassphrase        | The passphrase for privacy. | `-` |
| contextName              | The context name for SNMP v3 messages. | `-` |

### Reading Outputs

Outputs are referenced by name. A single device may have more than one instance
of an output type. A value of `-` in the table below indicates that there is no value
set for that field. The *custom* section describes outputs which this plugin defines
while the *built-in* section describes outputs this plugin uses which are [built-in to
the SDK](https://synse.readthedocs.io/en/latest/sdk/concepts/reading_outputs/#built-ins).

**Custom**

| Name    | Description                              | Unit  | Type    | Precision |
| ------- | ---------------------------------------- | :---: | ------- | :-------: |
| volt-ampere | A measure of power, in volt-amperes. | VA    | `power` | 3         |
| identity    | An output for SNMP identifiers.      | -     | `-`     | -         |

**Built-in**

| Name             | Description                       | Unit  | Type        | Precision |
| ---------------- | --------------------------------- | :---: | ----------- | :-------: |
| electric-current | A measure of current, in amperes. | A     | `current`   | 3         |
| frequency        | A measure of frequency, in hertz. | Hz    | `frequency` | 2         |
| watt             | A measure of power, in watts.     | W     | `watt`      | 3         |
| status           | A general measure of status.      | -     | `-`         | -         |

### Device Handlers

Device Handlers are referenced by name.

| Name      | Description                                    | Outputs            | Read  | Write | Bulk Read | Listen |
| --------- | ---------------------------------------------- | ------------------ | :---: | :---: | :-------: | :----: |
| current   | A handler for OIDs which report current.       | `electric-current` | ✓     | ✗     | ✗         | ✗      |
| frequency | A handler for OIDs which report frequency.     | `frequency`        | ✓     | ✗     | ✗         | ✗      |
| identity  | A handler for OIDs which report SNMP identity. | `identity`         | ✓     | ✗     | ✗         | ✗      |
| power     | A handler for OIDs which report power.         | `watt`             | ✓     | ✗     | ✗         | ✗      |
| status    | A handler for OIDs which report status.        | `status`           | ✓     | ✗     | ✗         | ✗      |

### Write Values

The SNMP plugin does not currently support writing to any of its device handlers.

## Supported MIBs

- [UPS-MIB][ups-mib-rfc]

## Compatibility

Below is a table describing the compatibility of plugin versions with Synse platform versions.

|             | Synse v2 | Synse v3 |
| ----------- | -------- | -------- |
| plugin v1.x | ✓        | ✗        |
| plugin v2.x | ✗        | ✓        |

## Troubleshooting

### Debugging

The plugin can be run in debug mode for additional logging. This is done by:

- Setting the `debug` option  to `true` in the plugin configuration YAML

  ```yaml
  debug: true
  ```

- Passing the `--debug` flag when running the binary/image

  ```
  docker run vaporio/snmp-plugin --debug
  ```

- Running the image with the `PLUGIN_DEBUG` environment variable set to `true`

  ```
  docker run -e PLUGIN_DEBUG=true vaporio/snmp-plugin
  ```

### Developing

A [development/debug Dockerfile](Dockerfile.dev) is provided in the project repository to enable
building image which may be useful when developing or debugging a plugin. Unlike the slim `scratch`-based
production image, the development image uses an ubuntu base, bringing with it all the standard command line
tools one would expect. To build a development image:

```
make docker-dev
```

The built image will be tagged using the format `dev-{COMMIT}`, where `COMMIT` is the short commit for
the repository at the time. This image is not published as part of the CI pipeline, but those with access
to the Docker Hub repo may publish manually.

## Contributing / Reporting

If you experience a bug, would like to ask a question, or request a feature, open a
[new issue](https://github.com/vapor-ware/synse-snmp-plugin/issues) and provide as much
context as possible. All contributions, questions, and feedback are welcomed and appreciated.

## License

The Synse SNMP Plugin is licensed under GPLv3. See [LICENSE](LICENSE) for more info.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-snmp-plugin.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fvapor-ware%2Fsynse-snmp-plugin?ref=badge_large)

[synse-server]: https://github.com/vapor-ware/synse-server
[synse-cli]: https://github.com/vapor-ware/synse-cli
[plugin-dockerhub]: https://hub.docker.com/r/vaporio/snmp-plugin
[plugin-release]: https://github.com/vapor-ware/synse-snmp-plugin/releases
[sdk-docs]: https://synse.readthedocs.io/en/latest/sdk/intro/
[dynamic-reg-example]: https://github.com/vapor-ware/synse-sdk/tree/master/examples/dynamic_registration
[ups-mib-rfc]: https://tools.ietf.org/html/rfc1628
