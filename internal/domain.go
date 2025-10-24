package internal

const AuthMethodBasic = "basic"

var authMethods = []string{AuthMethodBasic}

type Domain struct {
	AuthUser     string `mapstructure:"username"`
	AuthPassword string `mapstructure:"password"`
	AuthMethod   string `mapstructure:"auth_method"`
	Ip4Address   string `mapstructure:"ip4_address"`
	Ip6Address   string `mapstructure:"ip6_address"`
	Ip6HostId    string `mapstructure:"ip6_host_id"`
	Name         string `mapstructure:"domain"`
	RefreshUrl   string `mapstructure:"refresh_url"`
}

func NewDomain() *Domain {
	return &Domain{AuthMethod: AuthMethodBasic}
}

type Domains struct {
	List []Domain
}
