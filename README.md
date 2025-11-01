# Yet another Dynamic DNS client

This client aims to be simple and easy to use for updating ip addresses in dyn dns services. Plan is to make the integration into any system as easy as possible and to extend it with features step by step.

## Features
- Define static ipv4/ipv6 addresses per domain (will use it instead of identified wan ips)
- Define ipv6 host id/interface id per domain (will use prefix + host id instead of identified wan ipv6)
- Define refresh url with placeholders per domain
- Refresh URL templates support
- Refresh domains periodically
- Supports http basic authentication

## Refresh URL
A refresh URL contains all relevant information to update the configuration for a domain. It can consist of placeholders. The following placeholders are available:

| Placeholder  | Description                                                |
|--------------|------------------------------------------------------------|
| `<protocol>` | Protocol used to connect to the service (default: "https") |
| `<host>`     | Hostname of the service                                    |
| `<domain>`   | Name of your domain                                        |
| `<ip4>`      | Placeholder for the IPv4 address                           |
| `<ip6>`      | Placeholder for the IPv6 address                           |
| `<username>` | Username to authenticate at the service                    |
| `<password>` | Password to authenticate at the service                    |

## Usage with config file (`refresh`)
A config file has to be defined with the required data for refreshing one or more domain configurations.

Refresh url templates can be defined in the config file as well, which makes reusability very easy.

The config file must have the name `config.(toml|json|yaml)` and has to be placed in `/etc/yddns`, `~/.yddns` or in the same directory where the executable resides.

>[!NOTE]
> Supported config formats are `json`, `toml` and `yaml`.
>

>[!TIP]
> Check `config.toml.dist` for an example config.

>[!TIP]
> A different location for the config can be used via the flag `--config-file`.
### Help
```
Usage:
  yddns refresh [flags]

Flags:
  -c, --config-file string   Override default config using absolute file path
  -i, --interval int         Define refresh interval in seconds
  -p, --periodically         Refresh periodically

```
### Examples
#### Update domains one time
```shell
$ yddns refresh
```
#### Update domains periodically with an interval of 1800 seconds
```shell
$ yddns refresh -p -i 1800
```
#### Use different config location
```shell
$ yddns refresh -c /path/to/config.json
```

## Usage via cli (`refresh domain`)
### Help
```
Usage:
  yddns refresh domain [refresh-url | :template-name] [flags]

Flags:
      --domain string           Set name of the domain in the refresh URL [<domain>]
      --host string             Set host name of the service [<host>]
      --ip4-address string      Set IPv4 address instead determining via wan request [<ip4>]
      --ip6-address string      Set IPv6 address instead determining via wan request [<ip6>]
      --ip6-host-id string      Set IPv6 host id/interface id and use prefix + host id
      --password string         Set password used to authenticate [<password>]
      --protocol string         Set protocol in the refresh URL [<protocol>] (default "https")
      --request-method string   Set request method of the service (default "GET")
      --user-agent string       Set user agent in refresh requests
      --username string         Set username used to authenticate [<username>]
```
### Example
```shell
$ yddns refresh domain https://my-provider.tld/update?ip=<ip4>,<ip6>&some=value --username john --password topsecret --user-agent Mozilla
```
## Install from source
Clone the repo, build the command and create a config. That's basically it.
```shell
$ git clone https://github.com/drieschel/yddns.git && cd yddns
$ go build
$ vi config.toml
```

## TODO
- Refresh only on config changes
- Add more refresh url templates
- Support usernames and passwords in env variables
- ... more tbc ...