package internal

import (
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

const AUTH_METHOD_BASIC = "basic"
const IDENT_URL_IPV4 = "https://v4.ident.me"
const IDENT_URL_IPV6 = "https://v6.ident.me"

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
	Do(req *http.Request) (resp *http.Response, err error)
}

type Client struct {
	HttpClient HttpClient
}

func NewClient(httpClient HttpClient) *Client {
	return &Client{HttpClient: httpClient}
}

func (c *Client) Refresh(domain Domain) error {
	url, err := c.BuildRefreshUrl(domain)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if domain.AuthUser != "" && domain.AuthPassword != "" {
		switch domain.AuthMethod {
		case AUTH_METHOD_BASIC:
			request.SetBasicAuth(domain.AuthUser, domain.AuthPassword)
		}
	}

	resp, err := c.HttpClient.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode > 204 {
		responseBodyBytes, _ := io.ReadAll(resp.Body)
		errorString := strings.Trim(string(responseBodyBytes), " ")
		if errorString == "" {
			errorString = resp.Status
		}

		return errors.New(fmt.Sprintf("%s", errorString))
	}

	defer resp.Body.Close()

	return nil
}

func (c *Client) BuildRefreshUrl(domain Domain) (string, error) {
	replacements := map[string]string{"<username>": domain.AuthUser, "<password>": domain.AuthPassword, "<domain>": domain.Name}

	for _, ipVersion := range domain.IpVersions {
		key := "<ip" + strconv.Itoa(ipVersion) + "addr>"

		if strings.Contains(domain.RefreshUrl, key) {
			ip, err := c.DetermineWanIp(ipVersion)
			if err != nil {
				return "", err
			}

			replacements[key] = ip
		}
	}

	url := domain.RefreshUrl
	for search, replacement := range replacements {
		url = strings.Replace(url, search, replacement, -1)
	}

	return url, nil
}

func (c *Client) DetermineWanIp(ipVersion int) (string, error) {
	url, err := c.GetIdentUrl(ipVersion)
	if err != nil {
		return "", err
	}

	response, err := c.HttpClient.Get(url)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	ip, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

func (c *Client) IsIpVersion(version int) bool {
	return slices.Contains(slices.Collect(maps.Keys(c.GetIdentUrls())), version)
}

func (c *Client) GetIdentUrl(ipVersion int) (string, error) {
	if !c.IsIpVersion(ipVersion) {
		return "", errors.New(fmt.Sprintf("invalid ip version (%d)", ipVersion))
	}

	return c.GetIdentUrls()[ipVersion], nil
}

func (c *Client) GetIdentUrls() map[int]string {
	return map[int]string{4: IDENT_URL_IPV4, 6: IDENT_URL_IPV6}
}
