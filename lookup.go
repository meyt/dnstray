package main

import (
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

// CheckDNSHealth iterates through the servers and updates their latency.
func CheckDNSHealth(servers []*DNSServer) {
	for _, server := range servers {
		server.SetIsTesting(true)
		l1 := measureLatency(server.Dns1, config.TEST_DOMAINS)
		l2 := measureLatency(server.Dns2, config.TEST_DOMAINS)
		// If both fail, latency is -1.
		// If one fails, we take the healthy one.
		// If both work, we average them.
		if l1 == -1 && l2 == -1 {
			server.Latency = -1
		} else if l1 == -1 {
			server.Latency = l1
		} else if l2 == -1 {
			server.Latency = l1
		} else {
			server.Latency = (l1 + l2) / 2
		}
		server.SetIsTesting(false)
	}
}

func measureLatency(ip string, domains []string) int {
	if ip == "" {
		return -1
	}

	c := new(dns.Client)
	c.Timeout = 2 * time.Second
	var totalDuration time.Duration
	successCount := 0

	for _, domain := range domains {
		if !strings.HasSuffix(domain, ".") {
			domain += "."
		}

		m := new(dns.Msg)
		m.SetQuestion(domain, dns.TypeA)

		// Ensure the IP has a port
		addr := ip
		if _, _, err := net.SplitHostPort(ip); err != nil {
			addr = ip + ":53"
		}

		r, rtt, err := c.Exchange(m, addr)
		if err == nil && r.Rcode == dns.RcodeSuccess {
			totalDuration += rtt
			successCount++
		}
	}

	if successCount == 0 {
		return -1
	}

	return int(totalDuration.Milliseconds()) / successCount
}
