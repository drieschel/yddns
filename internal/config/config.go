package config

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/drieschel/yddns/internal/cache"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

const (
	AppName = "yddns"

	AuthMethodBasic  = "basic"
	AuthMethodBearer = "bearer"

	ProtocolHttp  = "http"
	ProtocolHttps = "https"

	RequestMethodGet  = "GET"
	RequestMethodPost = "POST"

	DefaultAppVersion      = "dev"
	DefaultAuthMethod      = AuthMethodBasic
	DefaultRequestMethod   = RequestMethodGet
	DefaultProtocol        = ProtocolHttps
	DefaultRefreshInterval = 60

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

	DirNameCache     = "cache"
	DirNameTemplates = "templates"

	RefreshUrlTemplatePrefix = ":"
)

var (
	Dirs                    = []string{fmt.Sprintf("/etc/%s", AppName), determineAppDir()}
	FilePath                = ""
	SupportedAuthMethods    = []string{AuthMethodBasic, AuthMethodBearer}
	SupportedProtocols      = []string{ProtocolHttp, ProtocolHttps}
	SupportedRequestMethods = []string{RequestMethodGet, RequestMethodPost}
)

type Config struct {
	AppDir          string
	AppVersion      string
	Domains         []*Domain            `mapstructure:"domains"`
	Templates       map[string]*Template `mapstructure:"templates"`
	RefreshInterval int                  `mapstructure:"refresh_interval"`
}

func NewConfig(appVersion string) *Config {
	c := &Config{
		AppDir:          determineAppDir(),
		AppVersion:      appVersion,
		Domains:         []*Domain{},
		Templates:       map[string]*Template{},
		RefreshInterval: DefaultRefreshInterval,
	}

	err := readFileTemplates(c)
	if err != nil {
		log.Fatal(err)
	}

	return c
}

func NewFileConfig(appVersion string) *Config {
	c := NewConfig(appVersion)

	err := readFileConfig(c)
	if err != nil {
		log.Fatal(err)
	}

	return c
}

func (c *Config) GetAppVersion() string {
	if c.AppVersion == "" {
		return DefaultAppVersion
	}

	return c.AppVersion
}

func (c *Config) GetDefaultUserAgent() string {
	return CreateDefaultUserAgent(c.GetAppVersion())
}

func (c *Config) GetTemplate(name string) (*Template, error) {
	if template, ok := c.Templates[name]; ok {
		return template, nil
	}

	return &Template{}, fmt.Errorf("template \"%s\" not found", name)
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

func (c *Config) PrepareAndGetDomains() ([]*Domain, error) {
	var domains []*Domain
	for _, d := range c.Domains {
		err := c.PrepareDomain(d)
		if err != nil {
			return domains, err
		}

		domains = append(domains, d)
	}

	return domains, nil
}

func (c *Config) CreateFileCache(expirySeconds int) cache.Cache {
	cacheDir := fmt.Sprintf("%s/%s", c.AppDir, DirNameCache)

	return cache.NewFileCache(cacheDir, expirySeconds)
}

func readFileTemplates(c *Config) error {
	fs := afero.NewOsFs()

	for i := len(Dirs) - 1; i >= 0; i-- {
		dir := Dirs[i]
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

			if _, exist := c.Templates[templateName]; exist {
				continue
			}

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

			c.Templates[templateName] = template
		}
	}

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

func CreateDefaultUserAgent(version string) string {
	return fmt.Sprintf("%s/%s", AppName, version)
}

func determineAppDir() string {
	execDir, err := filepath.Abs("./")
	if err != nil {
		log.Fatal(err)
	}

	parentDir, dirName := filepath.Split(execDir)
	if dirName == "bin" {
		return parentDir
	}

	return execDir
}
