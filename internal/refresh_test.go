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
		c.HttpClient = NewHttpClientMock(t, expectedIp, identUrl)
		actualIp := c.DetermineWanIp(ipVersion)

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
			assert.Equal(t, fmt.Sprintf("Invalid ip version (%d)", data.Version), err.Error())
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
	t                *testing.T
	ExpectedResponse string
	ExpectedUrl      string
}

func (c *HttpClientMock) Get(url string) (resp *http.Response, err error) {
	assert.Equal(c.t, c.ExpectedUrl, url)

	return &http.Response{
		Status:     "OK",
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBufferString(c.ExpectedResponse)),
	}, nil
}

func NewHttpClientMock(t *testing.T, expectedResponse string, expectedUrl string) *HttpClientMock {
	return &HttpClientMock{t: t, ExpectedResponse: expectedResponse, ExpectedUrl: expectedUrl}
}
