package TSDNS

import (
	"github.com/wolf-joe/ts-dns/DNSCaller"
	"github.com/wolf-joe/ts-dns/GFWList"
	"github.com/wolf-joe/ts-dns/Hosts"
	ipset "github.com/wolf-joe/ts-dns/IPSet"
	"github.com/wolf-joe/ts-dns/cache"
	"strings"
)

type Config struct {
	Cache        *cache.DNSCache
	Listen       string
	GFWChecker   *GFWList.DomainChecker
	HostsReaders []Hosts.Reader
	GroupMap     map[string]Group
}

type Group struct {
	Callers  []DNSCaller.Caller
	Matcher  *DomainMatcher
	IPSet    *ipset.IPSet
	IPSetTTL int
}

type DomainMatcher struct {
	checker *GFWList.DomainChecker
}

func (matcher *DomainMatcher) IsMatch(domain string) (match bool, ok bool) {
	return matcher.checker.IsBlocked(domain)
}

func NewDomainMatcher(rules []string) (matcher *DomainMatcher) {
	matcher = new(DomainMatcher)
	text := strings.Join(rules, "\n")
	matcher.checker, _ = GFWList.NewCheckerByStr(text, false)
	return
}
