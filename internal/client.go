package internal

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Client struct {
	AuthUser        string `mapstructure:"username"`
	AuthPassword    string `mapstructure:"password"`
	Domain          string `mapstructure:"domain"`
	IpVersions      []int  `mapstructure:"ip_versions"`
	PingUrlTemplate string `mapstructure:"url_template"`
}

func NewClient(authUser string, authPassword string, domain string, ipVersions []int, pingUrlTemplate string) *Client {
	return &Client{AuthUser: authUser, AuthPassword: authPassword, Domain: domain, IpVersions: ipVersions, PingUrlTemplate: pingUrlTemplate}
}

func (s *Client) Ping() error {
	pingUrl := s.BuildPingUrl()
	resp, err := http.Get(pingUrl)

	if err != nil {
		return err
	}

	if resp.StatusCode > 201 {
		return errors.New(fmt.Sprintf("Invalid status code: %d", resp.StatusCode))
	}

	defer resp.Body.Close()

	return nil
}

func (c *Client) BuildPingUrl() string {
	replacements := map[string]string{"<user>": c.AuthUser, "<password>": c.AuthPassword, "<Domain>": c.Domain}

	for ipVersion := range c.IpVersions {
		ipStr := strconv.Itoa(ipVersion)
		key := "<ip" + ipStr + "addr>"
		replacements[key] = DetermineIp(ipVersion)
	}

	url := c.PingUrlTemplate
	for search, replacement := range replacements {
		url = strings.Replace(url, search, replacement, -1)
	}

	return url
}
