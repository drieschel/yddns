# Yet another Dynamic DNS client

This client aims to be simple and easy to use for updating ip addresses in dyn dns services. Plan is to make the integration into any system as easy as possible and to extend it with features step by step.

### Features
- Define static ipv4/ipv6 addresses per domain (will use it instead of identified wan ips)
- Define ipv6 host id/interface id per domain (will use prefix + host id instead of identified wan ipv6)
- Define refresh url with placeholders per domain
  - Placeholders are `<domain> <username> <password> <ip4> <ip6>` (more may be added in future)
- Refresh domains periodically
- Supports http basic authentication

### TODO
- Refresh only on config changes
- Add refresh url templates for simple configuration
- Support usernames and passwords in env variables
- ... more soon ...

### Install from source
Clone the repo, build the command and create a config. That's basically it.
````shell
$ git clone https://github.com/drieschel/yddns.git && cd yddns
$ go build
$ vi config.toml
````
>[!NOTE]
> Supported config formats are `json`, `toml` and `yaml`.

>[!TIP]
> Check `config.toml.dist` for an example config.