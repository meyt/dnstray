package main

import (
	"log"
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
	initConfig(CONFIG_FILENAME, CONFIG)
	loadConfig(CONFIG_FILENAME)
	systray.Run(onReady, onExit)
}

func initConfig(filename string, text string) {
	if _, err := os.Stat(filename); err == nil {
		return
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
}

func loadConfig(filename string) {
	configData, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal([]byte(configData), &config)
	if err != nil {
		panic(err)
	}
}

func testDnsServers() {
	CheckDNSHealth(config.DNSServers)
	// for _, item := range dnsMenuItems {

	// }
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
		return
	}
	fastest := dnsServers[0]
	for _, item := range dnsServers[1:] {
		if item.Latency < fastest.Latency {
			fastest = item
		}
	}
	activateDnsServer(*fastest)
}

func loadMenuState() {
	// load current dns
	dnsServers, err := GetDNSServers()
	if err != nil {
		return
	}

	// update menu items
	for _, item := range config.DNSServers {
		dns1 := item.GetAddr1()
		dns2 := item.GetAddr2()
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
	item.SetIsApplying(true)
	SetDNS(item.Dns1, item.Dns2)
	time.Sleep(DNS_APPLY_WAIT) // wait to apply dns in linux
	item.SetIsApplying(false)
	loadMenuState()
}

func setupMenu() {

	for idx, server := range config.DNSServers {
		item := systray.AddMenuItem(
			getDNSMenuTitle(*server),
			server.Dns1+","+server.Dns2,
		)
		config.DNSServers[idx].Index = idx
		dnsMenuItems[idx] = item

		go func(item *systray.MenuItem) {
			for {
				<-item.ClickedCh
				activateDnsServer(*server)
			}
		}(item)
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
			SetDNS("", "")
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

func onExit() {}

func (item *DNSServer) GetAddr1() netip.Addr {
	return netip.MustParseAddr(item.Dns1)
}

func (item *DNSServer) GetAddr2() netip.Addr {
	return netip.MustParseAddr(item.Dns2)
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
