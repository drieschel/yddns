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

type Utils struct {
	HttpClient HttpClient
}

func NewUtils() *Utils {
	return &Utils{HttpClient: &http.Client{}}
}

func (u *Utils) CreatePingUrlTemplate(ipVersions []int) string {
	url := "https://<user>:<password>@dyndns.kasserver.com"

	query := make([]string, len(ipVersions))
	for i, ipVersion := range ipVersions {
		ipStr := strconv.Itoa(ipVersion)
		query[i] = "myip" + ipStr + "=<ip" + ipStr + "addr>"
	}

	return url + "?" + strings.Join(query, "&")
}

func (u *Utils) DetermineIp(ipVersion int) string {
	url, err := u.GetIdentUrl(ipVersion)
	if err != nil {
		log.Fatal(err)
	}

	response, err := u.HttpClient.Get(url)
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

func (u *Utils) IsIpVersion(version int) bool {
	return slices.Contains(slices.Collect(maps.Keys(u.GetIdentUrls())), version)
}

func (u *Utils) ValidateIpVersion(version int) {
	if !u.IsIpVersion(version) {
		log.Fatalf("Invalid ip version: %d", version)
	}
}

func (u *Utils) GetIdentUrl(ipVersion int) (string, error) {
	if !u.IsIpVersion(ipVersion) {
		return "", errors.New(fmt.Sprintf("Invalid ip version (%d)", ipVersion))
	}

	return u.GetIdentUrls()[ipVersion], nil
}

func (u *Utils) GetIdentUrls() map[int]string {
	return map[int]string{4: IDENT_URL_IPV4, 6: IDENT_URL_IPV6}
}
