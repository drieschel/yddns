# Yet another Dynamic DNS client

This dynamic dns client aims to be simple and easy to use. Plan is to make the integration into any system as easy as possible and to extend it with features step by step.

### Features
- Define static ipv4/ipv6 addresses per domain (will use it instead of identified wan ips)
- Define ipv6 host id/interface id per domain (will use prefix + host id instead of identified wan ipv6)
- Define refresh url with placeholders per domain
  - Placeholders are `<domain> <username> <password> <ip4> <ip6>` (more may be added in future)
- Refresh domains periodically
- Supports http basic authentication

### TODO
- Add systemd unit file for easy integration into linux systems
- Add refresh url templates for simple configuration
- Store usernames and passwords in env variables
- ... more soon ...

### Other
Check `config.toml.dist` for example configuration