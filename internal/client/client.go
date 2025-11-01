package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/drieschel/yddns/internal/config"
)

const (
	IdentUrlIpv4 = "https://v4.ident.me"
	IdentUrlIpv6 = "https://v6.ident.me"
)

type HttpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
	Get(url string) (resp *http.Response, err error)
}

type Client struct {
	httpClient HttpClient
	wanIp4     string
	wanIp6     string
}

func NewClient(httpClient HttpClient) *Client {
	return &Client{httpClient: httpClient}
}

func (c *Client) Refresh(domain *config.Domain) (string, error) {
	url, err := c.BuildRefreshUrl(domain)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest(domain.RequestMethod, url, nil)
	if err != nil {
		return "", err
	}

	switch domain.AuthMethod {
	case config.AuthMethodBasic:
		if domain.AuthUser != "" && domain.AuthPassword != "" {
			request.SetBasicAuth(domain.AuthUser, domain.AuthPassword)
		}

	case config.AuthMethodBearer:
		if domain.AuthPassword != "" {
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", domain.AuthPassword))
		}
	}

	request.Header.Set("User-Agent", domain.UserAgent)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", err
	}

	var responseBody []byte
	responseBody, err = io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	responseString := strings.Trim(string(responseBody), " ")

	if response.StatusCode > 204 {
		if responseString == "" {
			responseString = response.Status
		}

		return "", errors.New(fmt.Sprintf("%s", responseString))
	}

	defer response.Body.Close()

	return responseString, nil
}

func (c *Client) BuildRefreshUrl(domain *config.Domain) (string, error) {
	replacements, err := c.BuildReplacements(domain)
	if err != nil {
		return "", err
	}

	url := domain.RefreshUrl
	for search, replacement := range replacements {
		url = strings.Replace(url, search, replacement, -1)
	}

	return url, nil
}

func (c *Client) BuildReplacements(domain *config.Domain) (map[string]string, error) {
	replacements := NewReplacements()

	replacements.
		SetDefault(config.KeyProtocol, config.DefaultProtocol)

	replacements.
		Set(config.KeyUsername, domain.AuthUser).
		Set(config.KeyPassword, domain.AuthPassword).
		Set(config.KeyDomainName, domain.DomainName).
		Set(config.KeyHost, domain.Host).
		Set(config.KeyProtocol, domain.Protocol)

	ip4Key := createReplaceKey(config.KeyIp4)
	if strings.Contains(domain.RefreshUrl, ip4Key) {
		var err error
		ip4 := domain.Ip4Address

		if ip4 == "" {
			ip4, err = c.DetermineWanIp4()
			if err != nil {
				return map[string]string{}, err
			}
		}

		replacements.Set(config.KeyIp4, ip4)
	}

	ip6Key := createReplaceKey(config.KeyIp6)
	if strings.Contains(domain.RefreshUrl, ip6Key) {
		var err error
		ip6 := domain.Ip6Address

		if ip6 == "" {
			ip6, err = c.DetermineWanIp6()
			if err != nil {
				return map[string]string{}, err
			}

			if domain.Ip6HostId != "" {
				prefix := strings.Join(strings.Split(ip6, ":")[:4], ":")
				ip6 = fmt.Sprintf("%s:%s", prefix, domain.Ip6HostId)
			}
		}

		replacements.Set(config.KeyIp6, ip6)
	}

	return replacements.Build(), nil
}

func (c *Client) DetermineWanIp4() (string, error) {
	if c.wanIp4 == "" {
		response, err := c.httpClient.Get(IdentUrlIpv4)
		if err != nil {
			return "", err
		}

		defer response.Body.Close()

		ip, err := io.ReadAll(response.Body)
		if err != nil {
			return "", err
		}

		c.wanIp4 = string(ip)
	}

	return c.wanIp4, nil
}

func (c *Client) DetermineWanIp6() (string, error) {
	if c.wanIp6 == "" {
		response, err := c.httpClient.Get(IdentUrlIpv6)
		if err != nil {
			return "", err
		}

		defer response.Body.Close()

		ip, err := io.ReadAll(response.Body)
		if err != nil {
			return "", err
		}

		c.wanIp6 = string(ip)
	}

	return c.wanIp6, nil
}

func (c *Client) Clear() {
	c.wanIp4 = ""
	c.wanIp6 = ""
}

type Replacements struct {
	defaults map[string]string
	items    map[string]string
}

func NewReplacements() *Replacements {
	return &Replacements{defaults: map[string]string{}, items: map[string]string{}}
}

func (r *Replacements) Build() map[string]string {
	replacements := map[string]string{}
	for key, value := range r.items {
		if defaultValue, ok := r.defaults[key]; value == "" && ok {
			value = defaultValue
		}

		replacements[createReplaceKey(key)] = value
	}

	return replacements
}

func (r *Replacements) Set(key string, value string) *Replacements {
	r.items[key] = value

	return r
}

func (r *Replacements) SetDefault(key string, value string) *Replacements {
	r.defaults[key] = value

	return r
}

func createReplaceKey(key string) string {
	return fmt.Sprintf("<%s>", key)
}
