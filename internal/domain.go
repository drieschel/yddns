package internal

type Domain struct {
	AuthUser     string `mapstructure:"username"`
	AuthPassword string `mapstructure:"password"`
	Domain       string `mapstructure:"domain"`
	IpVersions   []int  `mapstructure:"ip_versions"`
	RefreshUrl   string `mapstructure:"refresh_url"`
}
