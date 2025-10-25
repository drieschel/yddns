package internal

import (
	"fmt"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
)

const (
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

	ConfigKeyDomains         = "domains"
	ConfigKeyRefreshInterval = "refresh_interval"

	DomainConfigKeyAuthMethod    = "auth_method"
	DomainConfigKeyDomain        = "domain"
	DomainConfigKeyHost          = "host"
	DomainConfigKeyIp4Address    = "ip4_address"
	DomainConfigKeyIp6Address    = "ip6_address"
	DomainConfigKeyIp6HostId     = "ip6_host_id"
	DomainConfigKeyPassword      = "password"
	DomainConfigKeyProtocol      = "protocol"
	DomainConfigKeyRefreshUrl    = "refresh_url"
	DomainConfigKeyRequestMethod = "request_method"
	DomainConfigKeyTemplate      = "template"
	DomainConfigKeyUserAgent     = "user_agent"
	DomainConfigKeyUsername      = "username"
)

var (
	DomainConfigDefaultValues = map[string]interface{}{DomainConfigKeyAuthMethod: DefaultAuthMethod, DomainConfigKeyRequestMethod: DefaultRequestMethod, DomainConfigKeyProtocol: DefaultProtocol, DomainConfigKeyUserAgent: CreateDefaultUserAgentByVersion("dev")}
	SupportedAuthMethods      = []string{AuthMethodBasic, AuthMethodBearer}
	SupportedProtocols        = []string{ProtocolHttp, ProtocolHttps}
	SupportedRequestMethods   = []string{RequestMethodGet, RequestMethodPost}
)

func CreateDefaultUserAgentByVersion(version string) string {
	return fmt.Sprintf("yddns/%s", version)
}

func CreateDomainConfigDefaultsHookFunc() mapstructure.DecodeHookFunc {
	return func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
		if to != reflect.TypeOf(Domain{}) {
			return data, nil
		}

		if from != reflect.TypeOf(map[string]interface{}{}) {
			return data, nil
		}

		for k, v := range DomainConfigDefaultValues {
			if _, ok := data.(map[string]interface{})[k]; !ok {
				data.(map[string]interface{})[k] = v
			}
		}

		return data, nil
	}
}
