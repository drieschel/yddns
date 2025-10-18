package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRefresh(t *testing.T) {
	for _, domain := range RefreshProvider() {
		client := Client{}
		request := createHttpRequest(domain)
		response := createHttpResponse("something")
		httpClientMock := NewMockHttpClient(t)
		httpClientMock.EXPECT().Do(request).Return(response, nil).Once()
		client.httpClient = httpClientMock

		err := client.Refresh(domain)

		assert.NoError(t, err)
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
	ip4 := uuid.New().String()
	ip6 := uuid.New().String()
	testData := BuildRefreshUrlProvider(ip4, ip6)

	httpClientMock := NewMockHttpClient(t)

	ip4Key := createReplaceKey("ip4addr")
	ip6Key := createReplaceKey("ip6addr")

	expectedIp4requests := 0
	expectedIp6requests := 0

	for _, data := range testData {
		if strings.Contains(data.Domain.RefreshUrl, ip4Key) {
			expectedIp4requests = 1
		}

		if strings.Contains(data.Domain.RefreshUrl, ip6Key) {
			expectedIp6requests = 1
		}
	}

	httpClientMock.EXPECT().Get(IDENT_URL_IPV4).Return(createHttpResponse(ip4), nil).Times(expectedIp4requests)
	httpClientMock.EXPECT().Get(IDENT_URL_IPV6).Return(createHttpResponse(ip6), nil).Times(expectedIp6requests)

	c.httpClient = httpClientMock

	for _, data := range testData {
		actualRefreshUrl, err := c.BuildRefreshUrl(data.Domain)

		assert.NoError(t, err)
		assert.Equal(t, data.ExpectedRefreshUrl, actualRefreshUrl)
	}
}

func BuildRefreshUrlProvider(ipv4 string, ipv6 string) []struct {
	Domain             Domain
	ExpectedRefreshUrl string
} {
	return []struct {
		Domain             Domain
		ExpectedRefreshUrl string
	}{
		{
			Domain:             Domain{AuthUser: "foo", AuthPassword: "bar", Name: "fooma.driescheldns.org", RefreshUrl: "https://fancy-dyn.dns?a=<username>&b=<password>&c=<domain>&e=<ip4addr>&f=<ip6addr>"},
			ExpectedRefreshUrl: fmt.Sprintf("https://fancy-dyn.dns?a=foo&b=bar&c=fooma.driescheldns.org&e=%s&f=%s", ipv4, ipv6),
		},
		{
			Domain:             Domain{AuthUser: "foo", AuthPassword: "bar", Name: "fooma.driescheldns.org", RefreshUrl: "https://fancy-dyn.dns?e=<ip4addr>"},
			ExpectedRefreshUrl: fmt.Sprintf("https://fancy-dyn.dns?e=%s", ipv4),
		},
		{
			Domain:             Domain{AuthUser: "foo", AuthPassword: "bar", Name: "fooma.driescheldns.org", RefreshUrl: "https://fancy-dyn.dns?a=<username>&b=<password>&f=<ip6addr>"},
			ExpectedRefreshUrl: fmt.Sprintf("https://fancy-dyn.dns?a=foo&b=bar&f=%s", ipv6),
		},
		{
			Domain:             Domain{AuthUser: "foo", AuthPassword: "bar", Name: "fooma.driescheldns.org", RefreshUrl: "https://fancy-dyn.dns?a=<username>&b=<password>&c=<domain>"},
			ExpectedRefreshUrl: fmt.Sprintf("https://fancy-dyn.dns?a=foo&b=bar&c=fooma.driescheldns.org"),
		},
	}
}

func TestDetermineWanIp4(t *testing.T) {
	c := Client{}

	expectedUrl := IDENT_URL_IPV4
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

func TestDetermineWanIp6(t *testing.T) {
	c := Client{}

	expectedUrl := IDENT_URL_IPV6
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
