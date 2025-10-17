package internal

import (
	"bytes"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDetermineWanIp(t *testing.T) {
	c := Client{}
	for ipVersion, expectedUrl := range c.GetIdentUrls() {
		expectedIp := uuid.New().String()

		httpClientMock := NewMockHttpClient(t)
		httpClientMock.EXPECT().Get(expectedUrl).Return(createHttpResponse(expectedIp), nil).Once()
		c.HttpClient = httpClientMock

		actualIp, err := c.DetermineWanIp(ipVersion)

		assert.NoError(t, err)
		assert.Equal(t, expectedIp, actualIp)
	}
}

func TestIsIpVersion(t *testing.T) {
	c := Client{}
	for _, data := range IpVersionsProvider() {
		assert.Equal(t, c.IsIpVersion(data.Version), data.IsValid)
	}
}

func TestGetIdentUrl(t *testing.T) {
	c := Client{}
	identUrls := c.GetIdentUrls()
	for _, data := range IpVersionsProvider() {
		actualUrl, err := c.GetIdentUrl(data.Version)
		if data.IsValid {
			assert.Nil(t, err)
			assert.Equal(t, identUrls[data.Version], actualUrl)
		} else {
			assert.Error(t, err)
			assert.Equal(t, fmt.Sprintf("invalid ip version (%d)", data.Version), err.Error())
		}
	}
}

func TestGetIdentUrls(t *testing.T) {
	c := Client{}
	identUrls := c.GetIdentUrls()
	assert.Equal(t, 2, len(identUrls))
	assert.Equal(t, "https://v4.ident.me", identUrls[4])
	assert.Equal(t, "https://v6.ident.me", identUrls[6])
}

func TestRefresh(t *testing.T) {
	for _, domain := range RefreshProvider() {
		client := Client{}
		request := createHttpRequest(domain)
		response := createHttpResponse("something")
		httpClientMock := NewMockHttpClient(t)
		httpClientMock.EXPECT().Do(request).Return(response, nil).Once()
		client.HttpClient = httpClientMock

		_ = client.Refresh(domain)
	}
}

func RefreshProvider() map[string]Domain {
	return map[string]Domain{
		"basic auth": Domain{RefreshUrl: uuid.New().String(), AuthUser: "test", AuthPassword: "test", AuthMethod: "basic"},
		"no auth":    Domain{RefreshUrl: uuid.New().String(), AuthUser: "", AuthPassword: "", AuthMethod: "basic"},
	}
}

func TestBuildRefreshUrl(t *testing.T) {
	c := Client{}
	ipv4 := uuid.New().String()
	ipv6 := uuid.New().String()
	testData := RefreshUrlProvider(ipv4, ipv6)

	httpClientMock := NewMockHttpClient(t)

	for _, data := range testData {
		for _, ipVersion := range data.Domain.IpVersions {
			url, _ := c.GetIdentUrl(ipVersion)

			response := ipv4
			if ipVersion == 6 {
				response = ipv6
			}

			//TODO: Extend with checks for ip address replacements in refresh url
			httpClientMock.EXPECT().Get(url).Return(createHttpResponse(response), nil).Once()
		}
	}

	c.HttpClient = httpClientMock

	for _, data := range testData {
		actualRefreshUrl, err := c.BuildRefreshUrl(data.Domain)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, data.ExpectedRefreshUrl, actualRefreshUrl)
	}
}

// TODO: Extend with checks for ip address replacements in refresh url
func RefreshUrlProvider(ipv4 string, ipv6 string) []struct {
	Domain             Domain
	ExpectedRefreshUrl string
} {
	return []struct {
		Domain             Domain
		ExpectedRefreshUrl string
	}{
		{
			Domain:             Domain{AuthUser: "foo", AuthPassword: "bar", Name: "fooma.driescheldns.org", IpVersions: []int{4, 6}, RefreshUrl: "https://fancy-dyn.dns?a=<username>&b=<password>&c=<domain>&e=<ip4addr>&f=<ip6addr>"},
			ExpectedRefreshUrl: fmt.Sprintf("https://fancy-dyn.dns?a=foo&b=bar&c=fooma.driescheldns.org&e=%s&f=%s", ipv4, ipv6),
		},
	}
}

func IpVersionsProvider() []struct {
	Version int
	IsValid bool
} {
	return []struct {
		Version int
		IsValid bool
	}{
		{Version: rand.IntN(3), IsValid: false},
		{Version: 4, IsValid: true},
		{Version: 5, IsValid: false},
		{Version: 6, IsValid: true},
		{Version: 7 + rand.IntN(992), IsValid: false},
	}
}

func createHttpRequest(domain Domain) *http.Request {
	req, _ := http.NewRequest("GET", domain.RefreshUrl, nil)

	if domain.AuthUser != "" && domain.AuthPassword != "" {
		req.SetBasicAuth(domain.AuthUser, domain.AuthPassword)
	}

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
