package internal

import "fmt"

const (
	authMethodBasic      = "basic"
	authMethodBearer     = "bearer"
	requestMethodGet     = "GET"
	requestMethodPost    = "POST"
	defaultAuthMethod    = authMethodBasic
	defaultRequestMethod = requestMethodGet
	defaultProtocol      = "https"
	replaceKeyDomainName = "domain"
	replaceKeyHost       = "host"
	replaceKeyPassword   = "password"
	replaceKeyProtocol   = "protocol"
	replaceKeyUsername   = "username"
)

var (
	supportedAuthMethods    = []string{authMethodBasic, authMethodBearer}
	supportedRequestMethods = []string{requestMethodGet, requestMethodPost}
)

type Domain struct {
	AuthMethod    string `mapstructure:"auth_method"`
	AuthUser      string `mapstructure:"username"`
	AuthPassword  string `mapstructure:"password"`
	DomainName    string `mapstructure:"domain"`
	Host          string `mapstructure:"host"`
	Ip4Address    string `mapstructure:"ip4_address"`
	Ip6Address    string `mapstructure:"ip6_address"`
	Ip6HostId     string `mapstructure:"ip6_host_id"`
	Protocol      string `mapstructure:"protocol"`
	RefreshUrl    string `mapstructure:"refresh_url"`
	RequestMethod string `mapstructure:"request_method"`
	Template      string `mapstructure:"template"`
	UserAgent     string `mapstructure:"user_agent"`
}

func NewDomain() *Domain {
	template := &Domain{AuthMethod: defaultAuthMethod, Protocol: defaultProtocol, RequestMethod: defaultRequestMethod}
	template.SetDefaultUserAgentByVersion("dev")

	return template
}

func (d *Domain) SetDefaultUserAgentByVersion(version string) *Domain {
	d.UserAgent = fmt.Sprintf("drieschel / yddns / %s", version)

	return d
}

type Domains struct {
	List []Domain
}
