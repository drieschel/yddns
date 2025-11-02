package config

import (
	"fmt"
	"reflect"
)

type Domain struct {
	Template     `mapstructure:",squash"`
	AuthUser     string `mapstructure:"username"`
	AuthPassword string `mapstructure:"password"`
	DomainName   string `mapstructure:"domain"`
	Ip4Address   string `mapstructure:"ip4_address"`
	Ip6Address   string `mapstructure:"ip6_address"`
	Ip6HostId    string `mapstructure:"ip6_host_id"`
}

func (d *Domain) GetTemplateName() (string, error) {
	if !d.RequiresTemplate() {
		return "", fmt.Errorf("refresh url \"%s\" is not prefixed as template (%s)", d.RefreshUrl, RefreshUrlTemplatePrefix)
	}

	return d.RefreshUrl[len(RefreshUrlTemplatePrefix):], nil
}

func (d *Domain) InitDefaultValues(appVersion string) {
	if d.AuthMethod == "" {
		d.AuthMethod = DefaultAuthMethod
	}

	if d.Protocol == "" {
		d.Protocol = DefaultProtocol
	}

	if d.RequestMethod == "" {
		d.RequestMethod = DefaultRequestMethod
	}

	if d.UserAgent == "" {
		d.UserAgent = CreateDefaultUserAgent(appVersion)
	}
}

func (d *Domain) MergeTemplate(t *Template) {
	tType := reflect.TypeOf(*t)
	tValue := reflect.ValueOf(*t)
	dValue := reflect.ValueOf(&d.Template).Elem()

	d.RefreshUrl = ""
	for i := 0; i < tType.NumField(); i++ {
		tFieldValue := tValue.Field(i)
		dFieldValue := dValue.Field(i)

		if dFieldValue.String() == "" && tFieldValue.String() != "" {
			dFieldValue.SetString(tFieldValue.String())
		}
	}
}

func (d *Domain) RequiresTemplate() bool {
	return len(d.RefreshUrl) > 0 && d.RefreshUrl[0:len(RefreshUrlTemplatePrefix)] == RefreshUrlTemplatePrefix
}

type Template struct {
	AuthMethod    string `json:"auth_method" mapstructure:"auth_method"`
	Host          string `json:"host" mapstructure:"host"`
	Protocol      string `json:"protocol" mapstructure:"protocol"`
	RefreshUrl    string `json:"refresh_url" mapstructure:"refresh_url"`
	RequestMethod string `json:"request_method" mapstructure:"request_method"`
	UserAgent     string `json:"user_agent" mapstructure:"user_agent"`
}
