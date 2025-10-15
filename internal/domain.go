package internal

type Domain struct {
	AuthUser     string `mapstructure:"username"`
	AuthPassword string `mapstructure:"password"`
	IpVersions   []int  `mapstructure:"ip_versions"`
	Name         string `mapstructure:"domain"`
	RefreshUrl   string `mapstructure:"refresh_url"`
}

func (d *Domain) InitDefaultValues() {
	if len(d.IpVersions) == 0 {
		d.IpVersions = []int{4}
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
