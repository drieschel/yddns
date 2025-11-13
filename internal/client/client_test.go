package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/drieschel/yddns/internal/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestClient_Refresh(t *testing.T) {
	for name, domain := range RefreshTable() {
		t.Run(name, func(t *testing.T) {
			expectedResponse := uuid.New().String()

			request := createHttpRequest(domain)
			response := createHttpResponse(expectedResponse)
			httpClientMock := NewMockHttpClient(t)
			httpClientMock.EXPECT().Do(request).Return(response, nil).Once()

			client := Client{httpClient: httpClientMock}

			actualResponse, err := client.Refresh(&domain)

			assert.NoError(t, err)
			assert.Equal(t, expectedResponse, actualResponse)
		})
	}
}

func RefreshTable() map[string]config.Domain {
	return map[string]config.Domain{
		"basic auth":     {Template: config.Template{RefreshUrl: uuid.New().String(), AuthMethod: "basic", UserAgent: "test"}, AuthUser: "test", AuthPassword: "pass"},
		"basic no auth ": {Template: config.Template{RefreshUrl: uuid.New().String(), AuthMethod: "basic", UserAgent: "test"}, AuthUser: "", AuthPassword: ""},
		"bearer auth":    {Template: config.Template{RefreshUrl: uuid.New().String(), AuthMethod: "bearer", UserAgent: "test"}, AuthUser: "not used", AuthPassword: "token"},
		"bearer no auth": {Template: config.Template{RefreshUrl: uuid.New().String(), AuthMethod: "bearer", UserAgent: "test"}, AuthUser: "not used", AuthPassword: ""},
		"no auth":        {Template: config.Template{RefreshUrl: uuid.New().String(), AuthMethod: "", UserAgent: "test"}, AuthUser: "test", AuthPassword: "foo"},
	}
}

func TestClient_BuildRefreshUrl(t *testing.T) {
	for _, data := range refreshUrlTable() {
		t.Run(data.name, func(t *testing.T) {
			httpClient := NewMockHttpClient(t)
			client := NewClient(httpClient)

			if data.wanIp4 != "" {
				httpClient.EXPECT().Get(IdentUrlIpv4).Return(createHttpResponse(data.wanIp4), nil).Once()
			}

			if data.wanIp6 != "" {
				httpClient.EXPECT().Get(IdentUrlIpv6).Return(createHttpResponse(data.wanIp6), nil).Once()
			}

			actualRefreshUrl, err := client.BuildRefreshUrl(&data.domain)

			assert.NoError(t, err)
			assert.Equal(t, data.expectedUrl, actualRefreshUrl)
		})
	}
}

func TestClient_BuildReplacements(t *testing.T) {
	for _, data := range refreshUrlTable() {
		t.Run(data.name, func(t *testing.T) {
			httpClient := NewMockHttpClient(t)
			client := NewClient(httpClient)

			if data.wanIp4 != "" {
				httpClient.EXPECT().Get(IdentUrlIpv4).Return(createHttpResponse(data.wanIp4), nil).Once()
			}

			if data.wanIp6 != "" {
				httpClient.EXPECT().Get(IdentUrlIpv6).Return(createHttpResponse(data.wanIp6), nil).Once()
			}

			replacements, err := client.BuildReplacements(&data.domain)

			assert.NoError(t, err)
			assert.Equal(t, data.expectedReplacements, replacements)
		})
	}
}

func refreshUrlTable() []struct {
	name                 string
	domain               config.Domain
	expectedReplacements map[string]string
	expectedUrl          string
	wanIp4               string
	wanIp6               string
} {
	return []struct {
		name                 string
		domain               config.Domain
		expectedReplacements map[string]string
		expectedUrl          string
		wanIp4               string
		wanIp6               string
	}{
		{
			name:                 "IPv4/IPv6 replacements WAN IPv4/IPv6 requests",
			domain:               config.Domain{DomainName: "yddns.drieschel.org", AuthUser: "foo", AuthPassword: "bar", Template: config.Template{RefreshUrl: "<protocol>://<host>?a=<username>&b=<password>&c=<domain>&e=<ip4>&f=<ip6>", Host: "fancy-dyn.dns", Protocol: "http"}},
			expectedReplacements: map[string]string{"<username>": "foo", "<password>": "bar", "<domain>": "yddns.drieschel.org", "<host>": "fancy-dyn.dns", "<protocol>": "http", "<ip4>": "125.148.255.41", "<ip6>": "e764:9ec5:88f3:94a9:ad4c:a7b4:4075:1ca7"},
			expectedUrl:          "http://fancy-dyn.dns?a=foo&b=bar&c=yddns.drieschel.org&e=125.148.255.41&f=e764:9ec5:88f3:94a9:ad4c:a7b4:4075:1ca7",
			wanIp4:               "125.148.255.41",
			wanIp6:               "e764:9ec5:88f3:94a9:ad4c:a7b4:4075:1ca7",
		},
		{
			name:                 "IPv4 replacement only WAN IPv4 request",
			domain:               config.Domain{Template: config.Template{RefreshUrl: "https://fancy-dyn.dns?e=<ip4>"}},
			expectedReplacements: map[string]string{"<username>": "", "<password>": "", "<domain>": "", "<host>": "", "<protocol>": "https", "<ip4>": "125.148.255.41"},
			expectedUrl:          "https://fancy-dyn.dns?e=125.148.255.41",
			wanIp4:               "125.148.255.41",
			wanIp6:               "",
		},
		{
			name:                 "IPv6 replacement only WAN IPv6 request",
			domain:               config.Domain{Template: config.Template{RefreshUrl: "https://fancy-dyn.dns?f=<ip6>"}},
			expectedReplacements: map[string]string{"<username>": "", "<password>": "", "<domain>": "", "<host>": "", "<protocol>": "https", "<ip6>": "e764:9ec5:88f3:94a9:ad4c:a7b4:4075:1ca7"},
			expectedUrl:          "https://fancy-dyn.dns?f=e764:9ec5:88f3:94a9:ad4c:a7b4:4075:1ca7",
			wanIp4:               "",
			wanIp6:               "e764:9ec5:88f3:94a9:ad4c:a7b4:4075:1ca7",
		},
		{
			name:                 "No IP replacements no WAN requests",
			domain:               config.Domain{Template: config.Template{RefreshUrl: "https://fancy-dyn.dns/something"}},
			expectedReplacements: map[string]string{"<username>": "", "<password>": "", "<domain>": "", "<host>": "", "<protocol>": "https"},
			expectedUrl:          fmt.Sprintf("https://fancy-dyn.dns/something"),
			wanIp4:               "",
			wanIp6:               "",
		},
		{
			name:                 "Static IPv4 and IPv6 no WAN requests IPv6 host id ignored",
			domain:               config.Domain{Ip4Address: "192.124.234.52", Ip6Address: "f724:a6ff:51dc:d827:5bbd:ce50:fa6a:d7e2", Ip6HostId: "a7cc:409a:e841:ea15", Template: config.Template{RefreshUrl: "https://fancy-dyn.dns?e=<ip4>&f=<ip6>"}},
			expectedReplacements: map[string]string{"<username>": "", "<password>": "", "<domain>": "", "<host>": "", "<protocol>": "https", "<ip4>": "192.124.234.52", "<ip6>": "f724:a6ff:51dc:d827:5bbd:ce50:fa6a:d7e2"},
			expectedUrl:          "https://fancy-dyn.dns?e=192.124.234.52&f=f724:a6ff:51dc:d827:5bbd:ce50:fa6a:d7e2",
			wanIp4:               "",
			wanIp6:               "",
		},
		{
			name:                 "IPv6 host id + WAN IPv6 request = WAN IPv6 prefix + host id",
			domain:               config.Domain{Ip4Address: "192.124.234.52", Ip6HostId: "a7cc:409a:e841:ea15", Template: config.Template{RefreshUrl: "https://fancy-dyn.dns?e=<ip4>&f=<ip6>"}},
			expectedReplacements: map[string]string{"<username>": "", "<password>": "", "<domain>": "", "<host>": "", "<protocol>": "https", "<ip4>": "192.124.234.52", "<ip6>": "d710:6c3b:b3c3:9f6b:a7cc:409a:e841:ea15"},
			expectedUrl:          "https://fancy-dyn.dns?e=192.124.234.52&f=d710:6c3b:b3c3:9f6b:a7cc:409a:e841:ea15",
			wanIp4:               "",
			wanIp6:               "d710:6c3b:b3c3:9f6b:a7cc:409a:e841:ea15",
		},
	}
}

func TestClient_DetermineWanIp4(t *testing.T) {
	c := Client{}

	expectedUrl := IdentUrlIpv4
	expectedIp := uuid.New().String()

	httpClientMock := NewMockHttpClient(t)
	httpClientMock.EXPECT().Get(expectedUrl).Return(createHttpResponse(expectedIp), nil).Once()
	c.httpClient = httpClientMock

	actualIp, err := c.DetermineWanIp4()

	assert.NoError(t, err)
	assert.Equal(t, expectedIp, actualIp)

	//determine ip cached
	actualIp, err = c.DetermineWanIp4()

	assert.NoError(t, err)
	assert.Equal(t, expectedIp, actualIp)
}

func TestClient_DetermineWanIp6(t *testing.T) {
	c := Client{}

	expectedUrl := IdentUrlIpv6
	expectedIp := uuid.New().String()

	httpClientMock := NewMockHttpClient(t)
	httpClientMock.EXPECT().Get(expectedUrl).Return(createHttpResponse(expectedIp), nil).Once()
	c.httpClient = httpClientMock

	actualIp, err := c.DetermineWanIp6()

	assert.NoError(t, err)
	assert.Equal(t, expectedIp, actualIp)

	//determine ip cached
	actualIp, err = c.DetermineWanIp6()

	assert.NoError(t, err)
	assert.Equal(t, expectedIp, actualIp)
}

func TestClient_Clear(t *testing.T) {
	client := Client{wanIp4: "foo", wanIp6: "bar"}
	client.Clear()
	assert.Empty(t, client.wanIp4)
	assert.Empty(t, client.wanIp6)
}

func createHttpRequest(domain config.Domain) *http.Request {
	req, _ := http.NewRequest("GET", domain.RefreshUrl, nil)

	if domain.AuthMethod == "basic" && domain.AuthUser != "" && domain.AuthPassword != "" {
		req.SetBasicAuth(domain.AuthUser, domain.AuthPassword)
	} else if domain.AuthMethod == "bearer" && domain.AuthPassword != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", domain.AuthPassword))
	}

	req.Header.Set("User-Agent", domain.UserAgent)

	return req
}

func createHttpResponse(responseBody string) *http.Response {
	return &http.Response{
		Status:     "OK",
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
	}
}
