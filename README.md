# Yet another Dynamic DNS client

This client aims to be simple and easy to use for updating ip addresses in dyn dns services. Plan is to make the integration into any system as easy as possible and to extend it with features step by step.

> [!CAUTION]
> The project is in alpha stadium. Please be aware that things like config structure and more may change until release.

## Features
- Define static ipv4/ipv6 addresses per domain (will use it instead of identified wan ips)
- Define ipv6 host id/interface id per domain (will use wan ipv6 prefix + host id instead of identified wan ipv6)
- Define refresh url with placeholders per domain
- Define refresh url templates in the config file
- Refresh domains periodically
- Supports authentication methods basic and bearer

## Refresh URL
A refresh URL contains all relevant information to update the configuration for a domain. It can consist of placeholders. The following placeholders are available:

| Placeholder  | Description                                                |
|--------------|------------------------------------------------------------|
| `<protocol>` | Protocol used to connect to the service (default: "https") |
| `<host>`     | Hostname of the service                                    |
| `<domain>`   | Name of your domain                                        |
| `<ip4>`      | Placeholder for the IPv4 address                           |
| `<ip6>`      | Placeholder for the IPv6 address                           |
| `<username>` | Username to authenticate on the service                    |
| `<password>` | Password to authenticate on the service                    |

## Domain config properties
For providing the best flexibility, the following configurable domain properties are available:

| Property       | Default value | Description                                                                                         |
|----------------|---------------|-----------------------------------------------------------------------------------------------------|
| refresh_url    | ""            | The refresh url or a template name. Template names must be prefixed with a colon (ie `":dyndns2"`). |
| username       | ""            | Can be used for basic authentication and in the refresh URL.                                        |
| password       | ""            | Can be used for basic and bearer authentication and in the refresh URL.                             |
| domain         | ""            | Can be used in the refresh URL, mostly used in combination with templates.                          |
| protocol       | "https"       | Can be used in the refresh URL, mostly used in combination with templates.                          | 
| host           | ""            | Can be used in the refresh URL, mostly used in combination with templates.                          |
| ip4_address    | ""            | A static IPv4 address can be provided.                                                              |
| ip6_address    | ""            | A static IPv6 address can be provided.                                                              |
| ip6_host_id    | ""            | A host id/interface id can be provided. Will be ignored in case `ip6_address` is defined.           |
| auth_method    | "basic"       | The authentication method for the service. Currently supported are "basic" and "bearer".            |
| request_method | "GET"         | Change the HTTP request method if necessary.                                                        |

## Usage with config file (`refresh`)
A config file has to be defined with the required data for refreshing one or more domain configurations.

Refresh url templates can be defined in the config file as well, which makes reusability very easy.

The config file must have the name `config.ext`, where `ext` represents the extension of a supported format. It has to be placed in `/etc/yddns`, `~/.yddns` or in the same directory where the executable resides.
>[!NOTE]
> Supported config formats are `json`, `toml` and `yaml`.

>[!TIP]
> Check `config.toml.example` for an example config with comments.
### Help
```
Usage:
  yddns refresh [flags]

Flags:
  -c, --config-file string     Override default config using absolute file path
  -p, --periodically           Refresh periodically
  -i, --refresh-interval int   Define refresh interval in seconds
```
### Examples
#### Update domain configurations one time
```shell
$ yddns refresh
```
#### Update domain configurations periodically with an interval of 1800 seconds
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
      --username string         Set username used for authentication [<username>]
      --password string         Set password used for authentication [<password>]
      --domain string           Set your dns domain [<domain>]
      --ip4-address string      Set IPv4 address instead determining via wan request [<ip4>]
      --ip6-address string      Set IPv6 address instead determining via wan request [<ip6>]
      --ip6-host-id string      Set IPv6 host id/interface id and use prefix + host id in the refresh url [<ip6>
      --host string             Set host name of the service in the refresh url [<host>]
      --protocol string         Set protocol in the refresh url [<protocol>] (default "https")
      --auth-method string      Set authentication method in refresh requests (default "basic")
      --request-method string   Set request method in refresh requests (default "GET")
      --user-agent string       Set user agent in refresh requests
      --cache-ttl int           Set relative domain configuration cache lifetime in seconds [0 is disabled] (default 600)
      --cache-max-ttl int       Set max domain configuration cache lifetime in seconds [0 is disabled] (default 86400)
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