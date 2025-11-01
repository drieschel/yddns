package config

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetTemplateName(t *testing.T) {
	expectedTemplateName := uuid.New().String()
	domain := Domain{Template: Template{RefreshUrl: fmt.Sprintf("%s%s", RefreshUrlTemplatePrefix, expectedTemplateName)}}
	actualTemplateName, err := domain.GetTemplateName()

	assert.NoError(t, err)
	assert.Equal(t, expectedTemplateName, actualTemplateName)
}

func TestGetTemplateNameWithoutTemplateName(t *testing.T) {
	templateName := uuid.New().String()
	domain := Domain{Template: Template{RefreshUrl: templateName}}
	_, err := domain.GetTemplateName()

	assert.Errorf(t, err, "refresh url \"%s\" is not prefixed as template (%s)", templateName, RefreshUrlTemplatePrefix)
}

func TestInitDefaultValues(t *testing.T) {
	tests := initDefaultValuesTables()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.givenDomain.InitDefaultValues("42.0")
			assert.Equal(t, test.expectedDomain, test.givenDomain)
		})
	}
}

func initDefaultValuesTables() []struct {
	name           string
	givenDomain    Domain
	expectedDomain Domain
} {
	return []struct {
		name           string
		givenDomain    Domain
		expectedDomain Domain
	}{
		{
			name:           "No default values set",
			givenDomain:    Domain{Template: Template{AuthMethod: "foo", Protocol: "bar", RequestMethod: "yes", UserAgent: "yddns/1.0.3"}},
			expectedDomain: Domain{Template: Template{AuthMethod: "foo", Protocol: "bar", RequestMethod: "yes", UserAgent: "yddns/1.0.3"}},
		},
		{
			name:           "All default values set",
			givenDomain:    Domain{Template: Template{AuthMethod: "", Protocol: "", RequestMethod: "", UserAgent: ""}},
			expectedDomain: Domain{Template: Template{AuthMethod: "basic", Protocol: "https", RequestMethod: "GET", UserAgent: "yddns/42.0"}},
		},
	}
}

func TestMergeTemplate(t *testing.T) {
	tests := mergeTemplateTables()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.NotEqual(t, test.expectedDomain, test.givenDomain)
			test.givenDomain.MergeTemplate(test.givenTemplate)
			assert.Equal(t, test.expectedDomain, test.givenDomain)
		})
	}
}

func mergeTemplateTables() []struct {
	name           string
	givenDomain    Domain
	givenTemplate  Template
	expectedDomain Domain
} {
	return []struct {
		name           string
		givenDomain    Domain
		givenTemplate  Template
		expectedDomain Domain
	}{
		{
			name: "domain template properties have higher prio except refresh_url",
			givenDomain: Domain{
				Template: Template{
					AuthMethod:    "secure",
					Host:          "foo",
					Protocol:      "bar",
					RefreshUrl:    "yalla",
					RequestMethod: "so",
					UserAgent:     "yddns/1.0.1",
				},
				AuthUser:     "jane",
				AuthPassword: "doe",
				DomainName:   "dome",
				Ip4Address:   "1.2.3.4",
				Ip6Address:   "::1",
				Ip6HostId:    "1:2:3:4",
			},
			givenTemplate: Template{
				AuthMethod:    "whatever",
				Host:          "one",
				Protocol:      "two",
				RefreshUrl:    "three",
				RequestMethod: "four",
				UserAgent:     "five",
			},
			expectedDomain: Domain{
				Template: Template{
					AuthMethod:    "secure",
					Host:          "foo",
					Protocol:      "bar",
					RefreshUrl:    "three",
					RequestMethod: "so",
					UserAgent:     "yddns/1.0.1",
				},
				AuthUser:     "jane",
				AuthPassword: "doe",
				DomainName:   "dome",
				Ip4Address:   "1.2.3.4",
				Ip6Address:   "::1",
				Ip6HostId:    "1:2:3:4",
			},
		},
		{
			name: "template properties are filling empty domain template properties",
			givenDomain: Domain{
				Template:     Template{},
				AuthUser:     "jane",
				AuthPassword: "doe",
				DomainName:   "dome",
				Ip4Address:   "1.2.3.4",
				Ip6Address:   "::1",
				Ip6HostId:    "1:2:3:4",
			},
			givenTemplate: Template{
				AuthMethod:    "whatever",
				Host:          "one",
				Protocol:      "two",
				RefreshUrl:    "three",
				RequestMethod: "four",
				UserAgent:     "five",
			},
			expectedDomain: Domain{
				Template: Template{
					AuthMethod:    "whatever",
					Host:          "one",
					Protocol:      "two",
					RefreshUrl:    "three",
					RequestMethod: "four",
					UserAgent:     "five",
				},
				AuthUser:     "jane",
				AuthPassword: "doe",
				DomainName:   "dome",
				Ip4Address:   "1.2.3.4",
				Ip6Address:   "::1",
				Ip6HostId:    "1:2:3:4",
			},
		},
	}
}

func TestRequiresTemplate(t *testing.T) {
	tests := requiresTemplateTable()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedValue, test.domain.RequiresTemplate())
		})
	}
}

func requiresTemplateTable() []struct {
	name          string
	domain        Domain
	expectedValue bool
} {
	return []struct {
		name          string
		domain        Domain
		expectedValue bool
	}{
		{
			name:          "Template name provided",
			domain:        Domain{Template: Template{RefreshUrl: ":template-it-is"}},
			expectedValue: true,
		},
		{
			name:          "Not a template name",
			domain:        Domain{Template: Template{RefreshUrl: "no-template-name"}},
			expectedValue: false,
		},
		{
			name:          "Empty not a template name either",
			domain:        Domain{Template: Template{RefreshUrl: ""}},
			expectedValue: false,
		},
	}
}

func TestUnmarshalDomain(t *testing.T) {
	for _, data := range unmarshalDomainTable() {
		t.Run("Unmarshal domain", func(t *testing.T) {
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
		})
	}
}

func unmarshalDomainTable() []domainTestData {
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
