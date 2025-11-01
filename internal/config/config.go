package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

const (
	AppName          = "yddns"
	AuthMethodBasic  = "basic"
	AuthMethodBearer = "bearer"

	ProtocolHttp  = "http"
	ProtocolHttps = "https"

	RequestMethodGet  = "GET"
	RequestMethodPost = "POST"

	DefaultAuthMethod           = AuthMethodBasic
	DefaultRequestMethod        = RequestMethodGet
	DefaultProtocol             = ProtocolHttps
	DefaultValueRefreshInterval = 600

	KeyDomains         = "domains"
	KeyRefreshInterval = "refresh_interval"

	KeyAuthMethod    = "auth_method"
	KeyDomainName    = "domain"
	KeyHost          = "host"
	KeyIp4Address    = "ip4_address"
	KeyIp6Address    = "ip6_address"
	KeyIp6HostId     = "ip6_host_id"
	KeyPassword      = "password"
	KeyProtocol      = "protocol"
	KeyRefreshUrl    = "refresh_url"
	KeyRequestMethod = "request_method"
	KeyUserAgent     = "user_agent"
	KeyUsername      = "username"

	KeyIp4 = "ip4"
	KeyIp6 = "ip6"

	DirNameTemplates = "templates"

	RefreshUrlTemplatePrefix = ":"
)

var (
	Dirs                    = []string{fmt.Sprintf("/etc/%s", AppName), fmt.Sprintf("%s/.%s", AppName, getHomeDir()), getExecDir()}
	DomainDefaultValues     = map[string]interface{}{KeyAuthMethod: DefaultAuthMethod, KeyRequestMethod: DefaultRequestMethod, KeyProtocol: DefaultProtocol, KeyUserAgent: CreateDefaultUserAgent("dev")}
	FilePath                = ""
	SupportedAuthMethods    = []string{AuthMethodBasic, AuthMethodBearer}
	SupportedProtocols      = []string{ProtocolHttp, ProtocolHttps}
	SupportedRequestMethods = []string{RequestMethodGet, RequestMethodPost}
)

type Config struct {
	AppVersion      string
	Domains         []Domain            `mapstructure:"domains"`
	Templates       map[string]Template `mapstructure:"templates"`
	RefreshInterval int                 `mapstructure:"refresh_interval"`
}

func NewConfig(appVersion string, fs afero.Fs) *Config {
	c := &Config{AppVersion: appVersion, Domains: []Domain{}}

	err := readTemplates(fs, c)
	if err != nil {
		log.Fatal(err)
	}

	return c
}

func NewFileConfig(appVersion string, fs afero.Fs) *Config {
	c := NewConfig(appVersion, fs)

	err := readFileConfig(c)
	if err != nil {
		log.Fatal(err)
	}

	return c
}

func (c *Config) GetAppVersion() string {
	if c.AppVersion == "" {
		return "dev"
	}

	return c.AppVersion
}

func (c *Config) GetDefaultUserAgent() string {
	return CreateDefaultUserAgent(c.GetAppVersion())
}

func (c *Config) GetTemplate(name string) (Template, error) {
	if template, ok := c.Templates[name]; ok {
		return template, nil
	}

	return Template{}, fmt.Errorf("template \"%s\" not found", name)
}

func CreateDefaultUserAgent(version string) string {
	return fmt.Sprintf("%s/%s", AppName, version)
}

func (c *Config) PrepareAndGetDomains() ([]Domain, error) {
	var domains []Domain
	for _, d := range c.Domains {
		err := c.PrepareDomain(&d)
		if err != nil {
			return domains, err
		}

		domains = append(domains, d)
	}

	return domains, nil
}

func (c *Config) PrepareDomain(d *Domain) error {
	if d.RequiresTemplate() {
		templateName, _ := d.GetTemplateName()
		template, err := c.GetTemplate(templateName)
		if err != nil {
			return err
		}

		d.MergeTemplate(template)
	}

	d.InitDefaultValues(c.GetAppVersion())

	return nil
}

func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return homeDir
}

func getExecDir() string {
	execDir, err := filepath.Abs("./")
	if err != nil {
		log.Fatal(err)
	}

	return execDir
}

func readTemplates(fs afero.Fs, c *Config) error {
	templates := map[string]Template{}
	for _, dir := range Dirs {
		templatesDir := filepath.Join(dir, DirNameTemplates)
		if exists, _ := afero.DirExists(fs, templatesDir); !exists {
			continue
		}

		jailedFs := afero.NewBasePathFs(fs, templatesDir)
		io := afero.NewIOFS(jailedFs)
		templateFiles, err := io.Glob("*.json")
		if err != nil {
			return err
		}

		for _, file := range templateFiles {
			templateName, _ := strings.CutSuffix(file, ".json")

			var data []byte
			data, err = io.ReadFile(file)
			if err != nil {
				return err
			}

			template := &Template{}
			err = json.Unmarshal(data, &template)
			if err != nil {
				return err
			}

			templates[templateName] = *template
		}

		break
	}

	c.Templates = templates

	return nil
}

func readFileConfig(c *Config) error {
	viper.SupportedExts = []string{"toml", "json", "yaml", "yml"}
	viper.SetConfigFile(FilePath)

	if FilePath == "" {
		viper.SetConfigName("config")
		for _, path := range Dirs {
			viper.AddConfigPath(path)
		}
	}

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(c)
	if err != nil {
		return err
	}

	return nil
}
