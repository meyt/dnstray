package main

var APP_WEBSITE = "https://github.com/meyt/dnstray"
var CONFIG_FILENAME = "config.toml"
var CONFIG = `
[[dns_servers]]
name = "Google"
dns1 = "8.8.8.8"
dns2 = "8.8.4.4"

[[dns_servers]]
name = "Cloudflare"
dns1 = "1.1.1.1"
dns2 = "1.0.0.1"

[[dns_servers]]
name = "OpenDNS"
dns1 = "208.67.222.222"
dns2 = "208.67.220.220"

[[dns_servers]]
name = "Quad9"
dns1 = "9.9.9.9"
dns2 = "149.112.112.112"

[[dns_servers]]
name = "AdGuard"
dns1 = "94.140.14.14"
dns2 = "94.140.15.15"

[[dns_servers]]
name = "IR - Recursive.dci.ir"
dns1 = "217.218.127.127"
dns2 = "217.218.155.155"

[[dns_servers]]
name = "IR - TCI"
dns1 = "5.200.200.200"
dns2 = "217.218.127.127"

[[dns_servers]]
name = "IR - Shecan.ir"
dns1 = "178.22.122.100"
dns2 = "185.51.200.2"

[[dns_servers]]
name = "IR - Begzar.ir"
dns1 = "185.55.226.26"
dns2 = "185.55.225.25"

[[dns_servers]]
name = "IR - Electrotm.org"
dns1 = "78.157.42.100"
dns2 = "78.157.42.101"

[[dns_servers]]
name = "IR - Radar.game"
dns1 = "10.202.10.10"
dns2 = "10.202.10.11"

[[dns_servers]]
name = "IR - DNSPro.ir"
dns1 = "87.107.110.109"
dns2 = "87.107.110.110"

[[dns_servers]]
name = "IR - Hostiran.net"
dns1 = "172.29.0.100"
dns2 = "172.29.2.100"

[[dns_servers]]
name = "IR - AsiaTech"
dns1 = "185.98.113.113"
dns2 = "185.98.114.114"

[[dns_servers]]
name = "IR - Shatel"
dns1 = "85.15.1.14"
dns2 = "85.15.1.15"

[[dns_servers]]
name = "IR - Pishgaman"
dns1 = "5.202.100.100"
dns2 = "5.202.100.101"

[[dns_servers]]
name = "IR - Mobinnet"
dns1 = "10.104.88.8"
dns2 = "8.8.8.8"

[[dns_servers]]
name = "IR - ParsOnline"
dns1 = "37.10.64.1"
dns2 = "37.10.65.1"

[[dns_servers]]
name = "IR - Sabanet"
dns1 = "89.40.90.100"
dns2 = "188.158.158.158"

[[dns_servers]]
name = "IR - Taknet"
dns1 = "185.47.48.122"
dns2 = "185.142.95.10"

[[dns_servers]]
name = "IR - ZiTel"
dns1 = "172.20.11.11"
dns2 = "172.20.11.12"
`
