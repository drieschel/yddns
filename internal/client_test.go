package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
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
	for _, data := range BuildRefreshUrlProvider() {
		httpClient := NewMockHttpClient(t)
		client := NewClient(httpClient)

		if data.WanIp4 != "" {
			httpClient.EXPECT().Get(IdentUrlIpv4).Return(createHttpResponse(data.WanIp4), nil).Once()
		}

		if data.WanIp6 != "" {
			httpClient.EXPECT().Get(IdentUrlIpv6).Return(createHttpResponse(data.WanIp6), nil).Once()
		}

		actualRefreshUrl, err := client.BuildRefreshUrl(data.Domain)

		assert.NoError(t, err)
		assert.Equal(t, data.ExpectedUrl, actualRefreshUrl)
	}
}

func BuildRefreshUrlProvider() []struct {
	Domain      Domain
	ExpectedUrl string
	WanIp4      string
	WanIp6      string
} {
	return []struct {
		Domain      Domain
		ExpectedUrl string
		WanIp4      string
		WanIp6      string
	}{
		{
			Domain:      Domain{Name: "yddns.drieschel.org", AuthUser: "foo", AuthPassword: "bar", RefreshUrl: "https://fancy-dyn.dns?a=<username>&b=<password>&c=<domain>&e=<ip4>&f=<ip6>"},
			ExpectedUrl: "https://fancy-dyn.dns?a=foo&b=bar&c=yddns.drieschel.org&e=125.148.255.41&f=e764:9ec5:88f3:94a9:ad4c:a7b4:4075:1ca7",
			WanIp4:      "125.148.255.41",
			WanIp6:      "e764:9ec5:88f3:94a9:ad4c:a7b4:4075:1ca7",
		},
		{
			Domain:      Domain{RefreshUrl: "https://fancy-dyn.dns?e=<ip4>"},
			ExpectedUrl: "https://fancy-dyn.dns?e=125.148.255.41",
			WanIp4:      "125.148.255.41",
			WanIp6:      "",
		},
		{
			Domain:      Domain{RefreshUrl: "https://fancy-dyn.dns?f=<ip6>"},
			ExpectedUrl: "https://fancy-dyn.dns?f=e764:9ec5:88f3:94a9:ad4c:a7b4:4075:1ca7",
			WanIp4:      "",
			WanIp6:      "e764:9ec5:88f3:94a9:ad4c:a7b4:4075:1ca7",
		},
		{
			Domain:      Domain{RefreshUrl: "https://fancy-dyn.dns/something"},
			ExpectedUrl: fmt.Sprintf("https://fancy-dyn.dns/something"),
			WanIp4:      "",
			WanIp6:      "",
		},
		{
			Domain:      Domain{Ip4Address: "192.124.234.52", Ip6Address: "f724:a6ff:51dc:d827:5bbd:ce50:fa6a:d7e2", Ip6HostId: "a7cc:409a:e841:ea15", RefreshUrl: "https://fancy-dyn.dns?e=<ip4>&f=<ip6>"},
			ExpectedUrl: "https://fancy-dyn.dns?e=192.124.234.52&f=f724:a6ff:51dc:d827:5bbd:ce50:fa6a:d7e2",
			WanIp4:      "",
			WanIp6:      "",
		},
		{
			Domain:      Domain{Ip4Address: "192.124.234.52", Ip6HostId: "a7cc:409a:e841:ea15", RefreshUrl: "https://fancy-dyn.dns?e=<ip4>&f=<ip6>"},
			ExpectedUrl: "https://fancy-dyn.dns?e=192.124.234.52&f=d710:6c3b:b3c3:9f6b:a7cc:409a:e841:ea15",
			WanIp4:      "",
			WanIp6:      "d710:6c3b:b3c3:9f6b:a7cc:409a:e841:ea15",
		},
	}
}

func TestDetermineWanIp4(t *testing.T) {
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

func TestDetermineWanIp6(t *testing.T) {
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
