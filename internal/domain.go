package internal

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

type Domains struct {
	List []Domain
}
