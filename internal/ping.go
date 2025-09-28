package internal

const IDENT_URL_IPV4 = "https://v4.ident.me/"
const IDENT_URL_IPV6 = "https://v6.ident.me/"

type Ping struct {
	AuthUser     string
	AuthPassword string
	Domain       string
	Protocols    []int
	Server       string
}

func NewDynDnsClient(authUser string, authPassword string, domain string, protocols []int, server string) *Ping {
	return &Ping{AuthUser: authUser, AuthPassword: authPassword, Domain: domain, Protocols: protocols, Server: server}
}

func (c *Ping) BuildRequestUrl() string {
	return ""
}

func validateProtocols(protocols []int) error {
	return nil
}
