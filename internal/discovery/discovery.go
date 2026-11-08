package discovery

import (
	"fmt"
	"net"
	"time"
)

// PeerInfo holds info about a discovered device on the local network
type PeerInfo struct {
	Name    string
	IP      string
	Port    int
	LastSeen time.Time
}

// GetLocalIPs returns all non-loopback local IP addresses
func GetLocalIPs() []string {
	var ips []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return ips
	}

	for _, iface := range ifaces {
		// skip down or loopback
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			// only IPv4 for simplicity
			if ip.To4() != nil {
				ips = append(ips, ip.String())
			}
		}
	}
	return ips
}

// GetPrimaryIP picks the best local IP (prefers 192.168.x.x or 10.x.x.x)
func GetPrimaryIP() string {
	ips := GetLocalIPs()
	if len(ips) == 0 {
		return "127.0.0.1"
	}

	// prefer private ranges
	for _, ip := range ips {
		parsed := net.ParseIP(ip)
		if parsed != nil && isPrivate(parsed) {
			return ip
		}
	}
	return ips[0]
}

func isPrivate(ip net.IP) bool {
	private := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}
	for _, cidr := range private {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// GetShareURL builds the full URL for sharing
func GetShareURL(ip string, port int) string {
	return fmt.Sprintf("http://%s:%d", ip, port)
}

// GetAllShareURLs returns URLs for all local IPs
func GetAllShareURLs(port int) []string {
	ips := GetLocalIPs()
	urls := make([]string, 0, len(ips))
	for _, ip := range ips {
		urls = append(urls, fmt.Sprintf("http://%s:%d", ip, port))
	}
	if len(urls) == 0 {
		urls = append(urls, fmt.Sprintf("http://127.0.0.1:%d", port))
	}
	return urls
}
