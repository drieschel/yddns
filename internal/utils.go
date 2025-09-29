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

func CreatePingUrlTemplate(ipVersions []int) string {
	url := "https://<user>:<password>@dyndns.kasserver.com"

	query := make([]string, len(ipVersions))
	for i, ipVersion := range ipVersions {
		ipStr := strconv.Itoa(ipVersion)
		query[i] = "myip" + ipStr + "=<ip" + ipStr + "addr>"
	}

	return url + "?" + strings.Join(query, "&")
}

func DetermineIp(ipVersion int) string {
	url, err := GetIdentUrl(ipVersion)
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.Get(url)
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

func IsIpVersion(version int) bool {
	return slices.Contains(slices.Collect(maps.Keys(GetIdentUrls())), version)
}

func ValidateIpVersion(version int) {
	if !IsIpVersion(version) {
		log.Fatalf("Invalid ip version: %d", version)
	}
}

func GetIdentUrl(ipVersion int) (string, error) {
	if !IsIpVersion(ipVersion) {
		return "", errors.New(fmt.Sprintf("invalid ip version (%d)", ipVersion))
	}

	return GetIdentUrls()[ipVersion], nil
}

func GetIdentUrls() map[int]string {
	return map[int]string{4: IDENT_URL_IPV4, 6: IDENT_URL_IPV6}
}
