package internal

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsIpVersion(t *testing.T) {
	for _, data := range IpVersionsProvider() {
		assert.Equal(t, IsIpVersion(data.Version), data.IsValid)
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
