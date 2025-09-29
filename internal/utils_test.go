package internal

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsIpVersion(t *testing.T) {
	for _, data := range IpVersionsProvider() {
		assert.Equal(t, IsIpVersion(data.Version), data.IsValid)
	}
}

func TestGetIdentUrl(t *testing.T) {
	identUrls := GetIdentUrls()
	for _, data := range IpVersionsProvider() {
		actualUrl, err := GetIdentUrl(data.Version)
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
	identUrls := GetIdentUrls()
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
