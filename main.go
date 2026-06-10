package main

import (
	"net/netip"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/meyt/dnstray/icon"

	"fyne.io/systray"
	"github.com/pelletier/go-toml"
	"github.com/skratchdot/open-golang/open"
)

type DNSServer struct {
	Index      int
	Name       string
	Dns1       string
	Dns2       string
	Latency    int
	IsApplying bool
	IsActive   bool
	IsTesting  bool
}

type Config struct {
	DNSServers   []*DNSServer `toml:"dns_servers"`
	TEST_DOMAINS []string     `toml:"test_domains"`
}

var config Config
var dnsMenuItems = map[int]*systray.MenuItem{} // index -> menu item

func main() {
	// Initialize logging first
	InitLogger()

	LogInfo("Starting dnstray...")

	// Catch panics and log them before exiting
	defer func() {
		if r := recover(); r != nil {
			LogError("PANIC: %v", r)
			LogError("Application crashed unexpectedly")
		}
	}()

	initConfig(CONFIG_FILENAME, CONFIG)
	loadConfig(CONFIG_FILENAME)
	systray.Run(onReady, onExit)
}

func initConfig(filename string, text string) {
	if _, err := os.Stat(filename); err == nil {
		return
	}

	LogInfo("Creating default config file: %s", filename)

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		LogError("failed creating config file: %s", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		LogError("failed writing to config file: %s", err)
		return
	}

	LogInfo("Default config file created successfully")
}

func loadConfig(filename string) {
	LogInfo("Loading config from: %s", filename)

	configData, err := os.ReadFile(filename)
	if err != nil {
		LogError("failed reading config file: %v", err)
		panic(err)
	}
	err = toml.Unmarshal([]byte(configData), &config)
	if err != nil {
		LogError("failed parsing config file: %v", err)
		panic(err)
	}

	LogInfo("Config loaded: %d DNS servers, %d test domains",
		len(config.DNSServers), len(config.TEST_DOMAINS))
}

func testDnsServers() {
	CheckDNSHealth(config.DNSServers)
}

func autoSelect() {
	CheckDNSHealth(config.DNSServers)
	dnsServers := []*DNSServer{}
	for _, item := range config.DNSServers {
		if item.Latency == 0 || item.Latency == -1 {
			continue
		}
		dnsServers = append(dnsServers, item)
	}
	if len(dnsServers) == 0 {
		LogWarn("No healthy DNS servers found for auto-select")
		return
	}
	fastest := dnsServers[0]
	for _, item := range dnsServers[1:] {
		if item.Latency < fastest.Latency {
			fastest = item
		}
	}
	LogInfo("Auto-selecting fastest DNS: %s (%dms)", fastest.Name, fastest.Latency)
	activateDnsServer(*fastest)
}

func loadMenuState() {
	// load current dns
	dnsServers, err := GetDNSServers()
	if err != nil {
		LogWarn("Failed to get current DNS servers: %v", err)
		return
	}

	// update menu items
	for _, item := range config.DNSServers {
		dns1, ok1 := safeParseAddr(item.Dns1)
		dns2, ok2 := safeParseAddr(item.Dns2)
		if !ok1 || !ok2 {
			continue
		}
		item.SetIsActive(slices.Contains(dnsServers, dns1) && slices.Contains(dnsServers, dns2))
	}
}

func getDNSMenuTitle(dnsServer DNSServer) string {
	r := dnsServer.Name
	if dnsServer.IsTesting {
		r += " [" + WAIT_MARK + "] "
	} else if dnsServer.Latency != 0 {
		r += " [" + strconv.Itoa(dnsServer.Latency) + "] "
	}
	if dnsServer.IsApplying {
		r += WAIT_MARK
	} else if dnsServer.IsActive {
		r += CHECK_MARK
	}
	return r
}

func activateDnsServer(item DNSServer) {
	LogInfo("Activating DNS server: %s (%s, %s)", item.Name, item.Dns1, item.Dns2)

	item.SetIsApplying(true)
	err := SetDNS(item.Dns1, item.Dns2)
	if err != nil {
		LogError("Failed to set DNS %s: %v", item.Name, err)
		item.SetIsApplying(false)
		return
	}
	time.Sleep(DNS_APPLY_WAIT) // wait to apply dns in linux
	item.SetIsApplying(false)
	loadMenuState()

	LogInfo("DNS server %s activated successfully", item.Name)
}

func setupMenu() {
	defer func() {
		if r := recover(); r != nil {
			LogError("PANIC in setupMenu: %v", r)
		}
	}()

	for idx, server := range config.DNSServers {
		item := systray.AddMenuItem(
			getDNSMenuTitle(*server),
			server.Dns1+","+server.Dns2,
		)
		config.DNSServers[idx].Index = idx
		dnsMenuItems[idx] = item

		go func(item *systray.MenuItem, server *DNSServer) {
			defer func() {
				if r := recover(); r != nil {
					LogError("PANIC in DNS menu click handler for %s: %v", server.Name, r)
				}
			}()
			for {
				<-item.ClickedCh
				activateDnsServer(*server)
			}
		}(item, server)
	}
	systray.AddSeparator()

	mClear := systray.AddMenuItem("Clear DNS", "Clear DNS settings")
	mAutoSelect := systray.AddMenuItem("Auto Select", "Test and select fastest")
	mTest := systray.AddMenuItem("Test", "Test DNS Servers")
	mAbout := systray.AddMenuItem("About", "About the app")
	mQuit := systray.AddMenuItem("Exit", "Quit the app")
	loadMenuState()

	for {
		select {
		case <-mClear.ClickedCh:
			LogInfo("Clearing DNS settings")
			err := SetDNS("", "")
			if err != nil {
				LogError("Failed to clear DNS: %v", err)
			}
			loadMenuState()
		case <-mAutoSelect.ClickedCh:
			mTest.Disable()
			mAutoSelect.Disable()
			mAutoSelect.SetTitle("Auto Select" + WAIT_MARK)
			autoSelect()
			mAutoSelect.SetTitle("Auto Select")
			mAutoSelect.Enable()
			mTest.Enable()
		case <-mTest.ClickedCh:
			mTest.Disable()
			mAutoSelect.Disable()
			mTest.SetTitle("Test" + WAIT_MARK)
			testDnsServers()
			mTest.SetTitle("Test")
			mAutoSelect.Enable()
			mTest.Enable()
		case <-mAbout.ClickedCh:
			open.Run(APP_WEBSITE)
		case <-mQuit.ClickedCh:
			LogInfo("User requested exit")
			systray.Quit()
			return
		}
	}
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("GO DNS Tray")
	systray.SetTooltip("Change system DNS")
	go setupMenu()
}

func onExit() {
	LogInfo("Application exiting")
}

// safeParseAddr safely parses an IP address string, returning the parsed
// address and true on success, or an empty address and false on failure.
// This replaces netip.MustParseAddr which panics on invalid input.
func safeParseAddr(s string) (netip.Addr, bool) {
	if s == "" {
		return netip.Addr{}, false
	}
	addr, err := netip.ParseAddr(s)
	if err != nil {
		LogWarn("Failed to parse DNS address '%s': %v", s, err)
		return netip.Addr{}, false
	}
	return addr, true
}

func (item *DNSServer) GetAddr1() (netip.Addr, bool) {
	return safeParseAddr(item.Dns1)
}

func (item *DNSServer) GetAddr2() (netip.Addr, bool) {
	return safeParseAddr(item.Dns2)
}

func (item *DNSServer) SetIsApplying(v bool) {
	item.IsApplying = v
	item.UpdateMenuTitle()
}

func (item *DNSServer) SetIsActive(v bool) {
	item.IsActive = v
	item.UpdateMenuTitle()
}

func (item *DNSServer) SetIsTesting(v bool) {
	item.IsTesting = v
	item.UpdateMenuTitle()
}

func (item *DNSServer) UpdateMenuTitle() {
	dnsMenuItems[item.Index].SetTitle(getDNSMenuTitle(*item))
}
