package internal

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type DynDnsClient struct {
	AuthUser     string
	AuthPassword string
	Domain       string
	IpVersions   []int
}

func NewDynDnsClient(authUser string, authPassword string, domain string, ipVersions []int) *DynDnsClient {
	return &DynDnsClient{AuthUser: authUser, AuthPassword: authPassword, Domain: domain, IpVersions: ipVersions}
}

func (s *DynDnsClient) Ping() error {
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

func (c *DynDnsClient) BuildPingUrl() string {
	replacements := map[string]string{"<user>": c.AuthUser, "<password>": c.AuthPassword, "<domain>": c.Domain}

	for ipVersion := range c.IpVersions {
		ipStr := strconv.Itoa(ipVersion)
		key := "<ip" + ipStr + "addr>"
		replacements[key] = DetermineIp(ipVersion)
	}

	url := CreatePingUrlTemplate(c.IpVersions)
	for search, replacement := range replacements {
		url = strings.Replace(url, search, replacement, -1)
	}

	return url
}
