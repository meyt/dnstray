package main

import (
	"log"
	"os"

	"github.com/getlantern/systray"
	"github.com/meyt/godnstray/icon"
	"github.com/pelletier/go-toml"
	"github.com/skratchdot/open-golang/open"
)

type DNSServer struct {
	Name string
	Dns1 string
	Dns2 string
}

type Config struct {
	DNSServers []DNSServer `toml:"dns_servers"`
}

var config Config

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

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("GO DNS Tray")
	systray.SetTooltip("Change system DNS")

	for _, server := range config.DNSServers {
		dns1 := server.Dns1
		dns2 := server.Dns2
		item := systray.AddMenuItem(server.Name, dns1)
		go func(item *systray.MenuItem) {
			for {
				<-item.ClickedCh
				SetDNS(dns1, dns2)
			}
		}(item)
	}

	systray.AddSeparator()

	mClear := systray.AddMenuItem("Clear DNS", "Clear DNS settings")
	mAbout := systray.AddMenuItem("About", "About the app")
	mQuit := systray.AddMenuItem("Exit", "Quit the app")

	for {
		select {
		case <-mClear.ClickedCh:
			SetDNS("", "")
		case <-mAbout.ClickedCh:
			open.Run(APP_WEBSITE)
		case <-mQuit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func onExit() {}
