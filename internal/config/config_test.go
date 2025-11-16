package config

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestConfig_NewFileConfig(t *testing.T) {
	thisDir, _ := filepath.Abs("./")
	Dirs = []string{fmt.Sprintf("%s/testdata", thisDir)}

	cfg := NewFileConfig("42")

	assert.Equal(t, 42, cfg.RefreshInterval)
	assert.Equal(t, 1, len(cfg.Domains))

	//All domain values are getting properly unmarshalled
	assert.Equal(t, "a1", cfg.Domains[0].AuthMethod)
	assert.Equal(t, "a2", cfg.Domains[0].Host)
	assert.Equal(t, "a3", cfg.Domains[0].Protocol)
	assert.Equal(t, "a4", cfg.Domains[0].DomainName)
	assert.Equal(t, "a5", cfg.Domains[0].Ip4Address)
	assert.Equal(t, "a6", cfg.Domains[0].Ip6Address)
	assert.Equal(t, "a7", cfg.Domains[0].Ip6HostId)
	assert.Equal(t, "a8", cfg.Domains[0].AuthPassword)
	assert.Equal(t, "a9", cfg.Domains[0].RefreshUrl)
	assert.Equal(t, "a10", cfg.Domains[0].AuthUser)
	assert.Equal(t, "a11", cfg.Domains[0].UserAgent)
	assert.Equal(t, "a12", cfg.Domains[0].RequestMethod)

	assert.Equal(t, 3, len(cfg.Templates))
	assert.Contains(t, cfg.Templates, "yddns")
	assert.Contains(t, cfg.Templates, "stat-dns")
	assert.Contains(t, cfg.Templates, "provider")

	//Template in config.ext wins over template in dedicated file
	//and all template values are getting properly unmarshalled
	assert.Equal(t, "y1", cfg.Templates["yddns"].AuthMethod)
	assert.Equal(t, "y2", cfg.Templates["yddns"].Host)
	assert.Equal(t, "y3", cfg.Templates["yddns"].Protocol)
	assert.Equal(t, "y4", cfg.Templates["yddns"].RefreshUrl)
	assert.Equal(t, "y5", cfg.Templates["yddns"].RequestMethod)
	assert.Equal(t, "y6", cfg.Templates["yddns"].UserAgent)

	//Dedicated file templates are getting properly unmarshalled as well
	assert.Equal(t, "p1", cfg.Templates["provider"].AuthMethod)
	assert.Equal(t, "p2", cfg.Templates["provider"].Host)
	assert.Equal(t, "p3", cfg.Templates["provider"].Protocol)
	assert.Equal(t, "p4", cfg.Templates["provider"].RefreshUrl)
	assert.Equal(t, "p5", cfg.Templates["provider"].RequestMethod)
	assert.Equal(t, "p6", cfg.Templates["provider"].UserAgent)
}

func TestConfig_GetAppVersion(t *testing.T) {
	cfg := &Config{}
	assert.Empty(t, cfg.AppVersion)
	assert.Equal(t, "dev", cfg.GetAppVersion())
	cfg.AppVersion = "foo"
	assert.Equal(t, "foo", cfg.GetAppVersion())
}

func TestConfig_GetDefaultUserAgent(t *testing.T) {
	version := uuid.New().String()
	cfg := &Config{AppVersion: version}
	assert.Equal(t, fmt.Sprintf("yddns/%s", version), cfg.GetDefaultUserAgent())
}

func TestConfig_GetTemplate(t *testing.T) {
	expectedTemplate := &Template{AuthMethod: "secure"}
	templates := map[string]*Template{"fancy": expectedTemplate}
	cfg := &Config{Templates: templates}

	actualTemplate, err := cfg.GetTemplate("fancy")

	assert.NoError(t, err)
	assert.Same(t, expectedTemplate, actualTemplate)
}

func TestConfig_GetTemplateNotFound(t *testing.T) {
	cfg := &Config{Templates: map[string]*Template{}}

	_, err := cfg.GetTemplate("404")
	assert.Errorf(t, err, "template \"%s\" not found", "404")
}

func TestCreateDefaultUserAgent(t *testing.T) {
	assert.Equal(t, "yddns/everything", CreateDefaultUserAgent("everything"))
}
