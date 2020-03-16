package outbound

import (
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/proxy"
	"net"
	"testing"
)

type DialerMock struct {
	proxy.Dialer
}

func (mock DialerMock) Dial(network, addr string) (c net.Conn, err error) {
	dialer := net.Dialer{}
	return dialer.Dial(network, addr)
}

var question = dns.Question{Name: "ip.cn.", Qtype: dns.TypeA}
var request, fakeRequest = &dns.Msg{}, &dns.Msg{}
var fakeQuest = dns.Question{Name: "ip.cn", Qtype: dns.TypeA}
var mockProxy = DialerMock{}
var fakeProxy, _ = proxy.SOCKS5("tcp", "unknown", nil, proxy.Direct)

func assertFail(t *testing.T, val *dns.Msg, err error) {
	assert.NotEqual(t, err, nil)
	assert.True(t, val == nil)
}
func assertSuccess(t *testing.T, val *dns.Msg, err error) {
	assert.Equal(t, err, nil)
	assert.True(t, len(val.Answer) > 0)
}

func TestDNSCaller(t *testing.T) {
	server := "1.1.1.1:53"
	request.SetQuestion(question.Name, question.Qtype)
	// 空初始化
	caller := NewDNSCaller("", "udp", nil)
	r, err := caller.Call(request)
	assertFail(t, r, err)
	// 无效dns地址
	caller = NewDNSCaller(server+"ne", "udp", nil)
	r, err = caller.Call(request)
	assertFail(t, r, err)
	// 正常请求
	caller = NewDNSCaller(server, "udp", nil)
	r, err = caller.Call(request)
	assertSuccess(t, r, err)
	// 无效请求
	fakeRequest.SetQuestion(fakeQuest.Name, fakeQuest.Qtype)
	r, err = caller.Call(fakeRequest)
	assertFail(t, r, err)
	// 无效代理
	caller = NewDNSCaller(server, "udp", fakeProxy)
	r, err = caller.Call(request)
	assertFail(t, r, err)
	// Mock代理
	caller = NewDNSCaller(server, "udp", mockProxy)
	r, err = caller.Call(request)
	assertSuccess(t, r, err)
}

func TestDohCaller(t *testing.T) {
	address, serverName := "1.0.0.1:853", "cloudflare-dns.com"
	// mock代理
	caller := NewDoTCaller(address, serverName, mockProxy)
	r, err := caller.Call(request)
	assertSuccess(t, r, err)
}

func TestDoHCaller(t *testing.T) {
	url := "https://cloudflare-dns.com/dns-query"
	// 无效服务器
	request.SetQuestion(question.Name, question.Qtype)
	caller := DoHCaller{Url: "https://not-exists.com/dns-query"}
	r, err := caller.Call(request)
	assertFail(t, r, err)
	// 无效路径
	caller = DoHCaller{Url: url + "/ne"}
	r, err = caller.Call(request)
	assertFail(t, r, err)
	// 正常请求
	caller = DoHCaller{Url: url}
	r, err = caller.Call(request)
	assertSuccess(t, r, err)
	// 无效请求
	fakeRequest.SetQuestion(fakeQuest.Name, fakeQuest.Qtype)
	r, err = caller.Call(fakeRequest)
	assertFail(t, r, err)
	// 无效代理
	caller = DoHCaller{Url: url, Dialer: fakeProxy}
	r, err = caller.Call(request)
	assertFail(t, r, err)
}
