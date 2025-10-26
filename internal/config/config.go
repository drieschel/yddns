package config

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

	KeyDomains         = "domains"
	KeyRefreshInterval = "refresh_interval"

	DomainKeyAuthMethod    = "auth_method"
	DomainKeyDomainName    = "domain"
	DomainKeyHost          = "host"
	DomainKeyIp4Address    = "ip4_address"
	DomainKeyIp6Address    = "ip6_address"
	DomainKeyIp6HostId     = "ip6_host_id"
	DomainKeyPassword      = "password"
	DomainKeyProtocol      = "protocol"
	DomainKeyRefreshUrl    = "refresh_url"
	DomainKeyRequestMethod = "request_method"
	DomainKeyUserAgent     = "user_agent"
	DomainKeyUsername      = "username"
)

var (
	DomainDefaultValues     = map[string]interface{}{DomainKeyAuthMethod: DefaultAuthMethod, DomainKeyRequestMethod: DefaultRequestMethod, DomainKeyProtocol: DefaultProtocol, DomainKeyUserAgent: CreateDefaultUserAgentByVersion("dev")}
	SupportedAuthMethods    = []string{AuthMethodBasic, AuthMethodBearer}
	SupportedProtocols      = []string{ProtocolHttp, ProtocolHttps}
	SupportedRequestMethods = []string{RequestMethodGet, RequestMethodPost}
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

		for k, v := range DomainDefaultValues {
			if _, ok := data.(map[string]interface{})[k]; !ok {
				data.(map[string]interface{})[k] = v
			}
		}

		return data, nil
	}
}
