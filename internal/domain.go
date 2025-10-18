package internal

var authMethods = []string{"basic"}

type Domain struct {
	AuthUser     string `mapstructure:"username"`
	AuthPassword string `mapstructure:"password"`
	AuthMethod   string `mapstructure:"auth_method"`
	IpVersions   []int  `mapstructure:"ip_versions"`
	Name         string `mapstructure:"domain"`
	RefreshUrl   string `mapstructure:"refresh_url"`
}

type Domains struct {
	List []Domain
}
