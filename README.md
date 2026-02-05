
# DNSTray

<p align="center">
  <img src="https://github.com/user-attachments/assets/feba7482-d373-45a8-8a32-b84acac19515" alt="DNSTray Screenshot" width="200">
</p>

<p align="center">
  <strong>Set DNS from your system tray with ease.</strong>
</p>

<p align="center">
  <a href="https://github.com/meyt/dnstray/releases"><img src="https://img.shields.io/github/v/release/meyt/dnstray?color=blue&label=Latest%20Release" alt="Latest Release"></a>
  <a href="https://github.com/meyt/dnstray/releases"><img src="https://img.shields.io/github/downloads/meyt/dnstray/total" alt="Total Downloads"></a>
  <a href="https://github.com/meyt/dnstray/blob/master/LICENSE"><img src="https://img.shields.io/github/license/meyt/dnstray" alt="License"></a>
</p>

## ✨ Features

- 🖱️ **Easy to use** – Control everything directly from the tray icon
- ⚙️ **Configurable** – Customize your DNS server list
- 🚀 **Auto-optimize** – Test and automatically select the fastest DNS server
- ✅ **Smart detection** – Automatically detects and marks the currently active DNS server

---

## 📥 Installation & Usage

### Windows 10/11

1. Download the appropriate binary for your architecture:

| Architecture | Download |
|-------------|----------|
| x64 | [dnstray.exe](https://github.com/meyt/dnstray/releases/latest/download/dnstray.exe) |
| ARM | [dnstray-arm.exe](https://github.com/meyt/dnstray/releases/latest/download/dnstray-arm.exe) |
| ARM64 | [dnstray-arm64.exe](https://github.com/meyt/dnstray/releases/latest/download/dnstray-arm64.exe) |

2. **Run as Administrator** (required for DNS configuration changes)

### Linux

```bash
# Download the latest release
wget -O dnstray https://github.com/meyt/dnstray/releases/latest/download/dnstray-linux-amd64

# Make it executable
chmod +x dnstray

# Run it
./dnstray
```

### macOS

> **Not yet...**

---

## ⚙️ Configuration

After the first run, DNSTray automatically creates a configuration file named `dnstray.toml` in the same directory.

### Customizing DNS Servers

1. Close DNSTray (right-click tray icon → Exit)
2. Open `dnstray.toml` in your favorite text editor
3. Modify the DNS server list
4. Save the file and restart DNSTray

### Example Configuration

```toml
test_domains = ["google.com", "bing.com"]

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
```

---

## 🛠️ Building from Source

```bash
git clone https://github.com/meyt/dnstray.git
cd dnstray
go mod tidy
go build
```

---

## 🤝 Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.
