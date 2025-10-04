package internal

import (
	"errors"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

const IDENT_URL_IPV4 = "https://v4.ident.me"
const IDENT_URL_IPV6 = "https://v6.ident.me"

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type Domain struct {
	AuthUser     string `mapstructure:"username"`
	AuthPassword string `mapstructure:"password"`
	Domain       string `mapstructure:"domain"`
	IpVersions   []int  `mapstructure:"ip_versions"`
	RefreshUrl   string `mapstructure:"refresh_url"`
}

type Client struct {
	Domains    []Domain
	HttpClient HttpClient
}

func NewClient(domains []Domain, httpClient HttpClient) *Client {
	return &Client{Domains: domains, HttpClient: httpClient}
}

func (c *Client) RefreshIp(domain Domain) error {
	refreshUrl := c.BuildRefreshUrl(domain)
	
	resp, err := c.HttpClient.Get(refreshUrl)
	if err != nil {
		return err
	}

	if resp.StatusCode > 201 {
		return errors.New(fmt.Sprintf("Invalid status code: %d", resp.StatusCode))
	}

	defer resp.Body.Close()

	return nil
}

func (c *Client) BuildRefreshUrl(domain Domain) string {
	replacements := map[string]string{"<username>": domain.AuthUser, "<password>": domain.AuthPassword, "<domain>": domain.Domain}

	for _, ipVersion := range domain.IpVersions {
		ipStr := strconv.Itoa(ipVersion)
		key := "<ip" + ipStr + "addr>"
		replacements[key] = c.DetermineWanIp(ipVersion)
	}

	url := domain.RefreshUrl
	for search, replacement := range replacements {
		url = strings.Replace(url, search, replacement, -1)
	}

	return url
}

func (c *Client) DetermineWanIp(ipVersion int) string {
	url, err := c.GetIdentUrl(ipVersion)
	if err != nil {
		log.Fatal(err)
	}

	response, err := c.HttpClient.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	ip, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(ip)
}

func (c *Client) IsIpVersion(version int) bool {
	return slices.Contains(slices.Collect(maps.Keys(c.GetIdentUrls())), version)
}

func (c *Client) ValidateIpVersion(version int) {
	if !c.IsIpVersion(version) {
		log.Fatalf("Invalid ip version: %d", version)
	}
}

func (c *Client) GetIdentUrl(ipVersion int) (string, error) {
	if !c.IsIpVersion(ipVersion) {
		return "", errors.New(fmt.Sprintf("Invalid ip version (%d)", ipVersion))
	}

	return c.GetIdentUrls()[ipVersion], nil
}

func (c *Client) GetIdentUrls() map[int]string {
	return map[int]string{4: IDENT_URL_IPV4, 6: IDENT_URL_IPV6}
}
