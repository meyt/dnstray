package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// SetDNS sets the DNS server addresses for all active network adapters.
// It uses PowerShell's Set-DnsClientServerAddress cmdlet which is available
// on Windows 8+ and Windows Server 2012+.
// Falls back to wmic for older systems if PowerShell fails.
func SetDNS(addr1 string, addr2 string) error {
	var isClearing = len(addr1) == 0 && len(addr2) == 0

	if isClearing {
		LogInfo("Clearing DNS settings for all active adapters")
		return setDNSClear()
	}

	servers := []string{}
	if addr1 != "" {
		servers = append(servers, addr1)
	}
	if addr2 != "" {
		servers = append(servers, addr2)
	}

	LogInfo("Setting DNS to %s for all active adapters", strings.Join(servers, ", "))
	return setDNSServers(servers)
}

func setDNSServers(servers []string) error {
	// Build PowerShell command using Set-DnsClientServerAddress
	// This works on Windows 8+ / Server 2012+
	quotedServers := make([]string, len(servers))
	for i, s := range servers {
		quotedServers[i] = fmt.Sprintf("'%s'", s)
	}
	serverList := strings.Join(quotedServers, ",")

	psScript := fmt.Sprintf(
		"$adapters = Get-NetAdapter | Where-Object { $_.Status -eq 'Up' }; "+
			"if ($adapters) { "+
			"  foreach ($adapter in $adapters) { "+
			"    Set-DnsClientServerAddress -InterfaceAlias $adapter.Name -ServerAddresses (%s) "+
			"  } "+
			"} else { "+
			"  Write-Error 'No active network adapters found' "+
			"}",
		serverList,
	)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		LogError("PowerShell Set-DnsClientServerAddress failed: %v, output: %s", err, string(output))

		// Fallback to wmic for older Windows systems
		LogWarn("Attempting fallback to wmic...")
		return setDNSWmic(servers)
	}

	if strings.Contains(string(output), "Write-Error") || strings.Contains(string(output), "error") {
		LogWarn("PowerShell reported issue: %s", string(output))
	}

	LogInfo("DNS set successfully via PowerShell")
	return nil
}

func setDNSClear() error {
	psScript :=
		"$adapters = Get-NetAdapter | Where-Object { $_.Status -eq 'Up' }; " +
			"if ($adapters) { " +
			"  foreach ($adapter in $adapters) { " +
			"    Set-DnsClientServerAddress -InterfaceAlias $adapter.Name -ResetServerAddresses " +
			"  } " +
			"}"

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		LogError("PowerShell clear DNS failed: %v, output: %s", err, string(output))

		// Fallback to wmic for older Windows systems
		LogWarn("Attempting fallback to wmic...")
		return setDNSWmicClear()
	}

	LogInfo("DNS cleared successfully via PowerShell")
	return nil
}

// setDNSWmic is a fallback for older Windows systems that don't have
// the Set-DnsClientServerAddress PowerShell cmdlet.
func setDNSWmic(servers []string) error {
	var addr string
	if len(servers) >= 2 {
		addr = fmt.Sprintf("(%s,%s)", servers[0], servers[1])
	} else if len(servers) == 1 {
		addr = fmt.Sprintf("(%s)", servers[0])
	}

	cmd := exec.Command("cmd", "/C", "wmic", "nicconfig", "where", "(IPEnabled=TRUE)", "call", "SetDNSServerSearchOrder", addr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		LogError("wmic SetDNSServerSearchOrder failed: %v, output: %s", err, string(output))
		return fmt.Errorf("failed to set DNS via both PowerShell and wmic: %w", err)
	}

	LogInfo("DNS set successfully via wmic fallback")
	return nil
}

func setDNSWmicClear() error {
	cmd := exec.Command("cmd", "/C", "wmic", "nicconfig", "where", "(IPEnabled=TRUE)", "call", "SetDNSServerSearchOrder", "()")
	output, err := cmd.CombinedOutput()
	if err != nil {
		LogError("wmic clear DNS failed: %v, output: %s", err, string(output))
		return fmt.Errorf("failed to clear DNS via both PowerShell and wmic: %w", err)
	}

	LogInfo("DNS cleared successfully via wmic fallback")
	return nil
}
