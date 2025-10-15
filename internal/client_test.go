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
	for ipVersion, identUrl := range c.GetIdentUrls() {
		expectedIp := uuid.New().String()
		c.HttpClient = NewHttpClientMockSingleRequest(t, expectedIp, identUrl)
		actualIp, err := c.DetermineWanIp(ipVersion)
		if err != nil {
			t.Fatal(err)
		}

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
	expectedRefreshUrl := uuid.New().String()
	domain := Domain{RefreshUrl: expectedRefreshUrl}
	client := Client{}
	client.HttpClient = NewHttpClientMock(t, []HttpClientMockData{*NewHttpClientMockData("doesnt matter", expectedRefreshUrl)})

	client.Refresh(domain)
}

func TestBuildRefreshUrl(t *testing.T) {
	c := Client{}
	ipv4 := uuid.New().String()
	ipv6 := uuid.New().String()
	testData := RefreshUrlProvider(ipv4, ipv6)
	var expectedHttpClientData []HttpClientMockData
	for _, data := range testData {
		for _, ipVersion := range data.Domain.IpVersions {
			url, err := c.GetIdentUrl(ipVersion)
			if err != nil {
				t.Fatal(err)
			}

			response := ipv4
			if ipVersion == 6 {
				response = ipv6
			}

			expectedHttpClientData = append(expectedHttpClientData, *NewHttpClientMockData(response, url))
		}
	}

	c.HttpClient = NewHttpClientMock(t, expectedHttpClientData)

	for _, data := range testData {
		actualRefreshUrl, err := c.BuildRefreshUrl(data.Domain)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, data.ExpectedRefreshUrl, actualRefreshUrl)
	}
}

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

type HttpClientMock struct {
	t           *testing.T
	Call        int
	alwaysFirst bool
	expected    []HttpClientMockData
}

func NewHttpClientMock(t *testing.T, expected []HttpClientMockData) *HttpClientMock {
	return &HttpClientMock{t: t, expected: expected, Call: 1, alwaysFirst: false}
}

func NewHttpClientMockSingleRequest(t *testing.T, expectedResponse string, expectedUrl string) *HttpClientMock {
	return NewHttpClientMock(t, []HttpClientMockData{{Response: expectedResponse, Url: expectedUrl}})
}

func NewHttpClientMockAlwaysFirst(t *testing.T, expectedResponse string, expectedUrl string) *HttpClientMock {
	client := NewHttpClientMock(t, []HttpClientMockData{{Response: expectedResponse, Url: expectedUrl}})
	client.alwaysFirst = true

	return client
}

func (c *HttpClientMock) Get(url string) (resp *http.Response, err error) {
	i := c.Call - 1
	if c.alwaysFirst {
		i = 0
	}

	if len(c.expected) <= i {
		c.t.Fatalf("No mock data exists for request #%d", c.Call)
	}

	expected := c.expected[i]

	assert.Equal(c.t, expected.Url, url)

	response := &http.Response{
		Status:     "OK",
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBufferString(expected.Response)),
	}

	c.Call++

	return response, nil
}

type HttpClientMockData struct {
	Response string
	Url      string
}

func NewHttpClientMockData(response string, url string) *HttpClientMockData {
	return &HttpClientMockData{Response: response, Url: url}
}
