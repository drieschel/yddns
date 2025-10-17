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

func (d *Domain) InitDefaultValues() {
	if len(d.IpVersions) == 0 {
		d.IpVersions = []int{4}
	}

	if d.AuthMethod == "" {
		d.AuthMethod = "basic"
	}
}

type Domains struct {
	List []Domain
}

func (d *Domains) InitDefaultValues() {
	for i := range d.List {
		d.List[i].InitDefaultValues()
	}
}
