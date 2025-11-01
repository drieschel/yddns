package config

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalDomain(t *testing.T) {
	for _, data := range unmarshalDomainProvider() {
		var (
			err          error
			actualResult Domain
		)

		viper.SetConfigType(data.givenType)
		reader := strings.NewReader(data.givenConfig)
		err = viper.ReadConfig(reader)

		assert.NoError(t, err)

		err = viper.Unmarshal(&actualResult)

		assert.NoError(t, err)
		assert.Equal(t, data.expectedResult, actualResult)
	}
}

func unmarshalDomainProvider() []domainTestData {
	return []domainTestData{
		{
			givenType: "json",
			givenConfig: `{
							 "auth_method": "auth",
     						 "domain": "domain1.tld",
							 "host": "horst1.tld",
							 "protocol": "foo",
      						 "username": "john",
      						 "password": "doe",
      						 "ip4_address": "127.0.0.1",
							 "ip6_address": "::1",
  							 "ip6_host_id": "0000:0000:0000:0001",
      						 "refresh_url": "abcde",
      						 "template": "somewhere",
      						 "request_method": "bar",
      						 "user_agent": "not-mozilla"
    					  }`,
			expectedResult: Domain{
				Template: Template{
					AuthMethod:    "auth",
					Host:          "horst1.tld",
					Protocol:      "foo",
					RefreshUrl:    "abcde",
					RequestMethod: "bar",
					UserAgent:     "not-mozilla",
				},
				AuthUser:     "john",
				AuthPassword: "doe",
				DomainName:   "domain1.tld",
				Ip4Address:   "127.0.0.1",
				Ip6Address:   "::1",
				Ip6HostId:    "0000:0000:0000:0001",
			},
		},
	}
}

type domainTestData struct {
	givenType      string
	givenConfig    string
	expectedResult Domain
}
