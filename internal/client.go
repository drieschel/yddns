package internal

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func (c *Client) Refresh(domain *Domain) error {
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
		case "", authMethodBasic:
			request.SetBasicAuth(domain.AuthUser, domain.AuthPassword)
		}
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode > 204 {
		var responseBodyBytes []byte
		responseBodyBytes, err = io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		errorString := strings.Trim(string(responseBodyBytes), " ")
		if errorString == "" {
			errorString = response.Status
		}

		return errors.New(fmt.Sprintf("%s", errorString))
	}

	defer response.Body.Close()

	return nil
}

func (c *Client) BuildRefreshUrl(domain *Domain) (string, error) {
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

func (c *Client) BuildReplacements(domain *Domain) (map[string]string, error) {
	replacements := map[string]string{}
	replacements[createReplaceKey(replaceKeyUsername)] = domain.AuthUser
	replacements[createReplaceKey(replaceKeyPassword)] = domain.AuthPassword
	replacements[createReplaceKey(replaceKeyDomainName)] = domain.DomainName
	replacements[createReplaceKey(replaceKeyHost)] = domain.Host
	replacements[createReplaceKey(replaceKeyProtocol)] = domain.Protocol

	ip4Key := createReplaceKey("ip4")
	if strings.Contains(domain.RefreshUrl, ip4Key) {
		var err error
		ip4 := domain.Ip4Address

		if ip4 == "" {
			ip4, err = c.DetermineWanIp4()
			if err != nil {
				return replacements, err
			}
		}

		replacements[ip4Key] = ip4
	}

	ip6Key := createReplaceKey("ip6")
	if strings.Contains(domain.RefreshUrl, ip6Key) {
		var err error
		ip6 := domain.Ip6Address

		if ip6 == "" {
			ip6, err = c.DetermineWanIp6()
			if err != nil {
				return replacements, err
			}

			if domain.Ip6HostId != "" {
				prefix := strings.Join(strings.Split(ip6, ":")[:4], ":")
				ip6 = fmt.Sprintf("%s:%s", prefix, domain.Ip6HostId)
			}
		}

		replacements[ip6Key] = ip6
	}

	return replacements, nil
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

func createReplaceKey(name string) string {
	return fmt.Sprintf("<%s>", name)
}
