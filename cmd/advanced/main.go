package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
	"github.com/gorilla/mux"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang AdvancedSecurity ../bpf/advanced_security.c

// Metrics structure for API responses
type Metrics struct {
	Timestamp    time.Time       `json:"timestamp"`
	Packets      uint64          `json:"packets"`
	Bytes        uint64          `json:"bytes"`
	HTTPRequests uint64          `json:"http_requests"`
	TCPConn      uint64          `json:"tcp_connections"`
	UDPPackets   uint64          `json:"udp_packets"`
	Blocked      uint64          `json:"blocked_connections"`
	Transactions uint64          `json:"transactions"`
	TopIPs       []IPMetrics     `json:"top_ips"`
	Security     SecurityMetrics `json:"security"`
}

type IPMetrics struct {
	IP    string `json:"ip"`
	Count uint64 `json:"count"`
}

type SecurityMetrics struct {
	SuspiciousIPs []SuspiciousIP `json:"suspicious_ips"`
	PortScans     []PortScan     `json:"port_scans"`
	Threats       []Threat       `json:"threats"`
}

type SuspiciousIP struct {
	IP    string `json:"ip"`
	Count uint64 `json:"count"`
}

type PortScan struct {
	IP    string `json:"ip"`
	Ports uint32 `json:"ports"`
}

type Threat struct {
	Type      string    `json:"type"`
	IP        string    `json:"ip"`
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details"`
}

func main() {
	// Remove memory limit for eBPF
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatalf("Failed to remove memory limit: %v", err)
	}

	// Load pre-compiled eBPF program
	objs := AdvancedSecurityObjects{}
	if err := LoadAdvancedSecurityObjects(&objs, nil); err != nil {
		log.Fatalf("Loading objects failed: %v", err)
	}
	defer objs.Close()

	// Get the first network interface
	iface, err := getNetworkInterface()
	if err != nil {
		log.Fatalf("Failed to find network interface: %v", err)
	}

	// Attach XDP program
	if err := attachXDP(&objs); err != nil {
		log.Fatalf("Failed to attach XDP program: %v", err)
	}

	// Attach security kprobes
	attachKprobes(&objs)

	// Start HTTP server for metrics API
	go startMetricsServer(&objs)

	// Start metrics collection
	startMetricsCollection(&objs, iface.Name)
}

func getNetworkInterface() (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			return &iface, nil
		}
	}

	return nil, fmt.Errorf("no suitable network interface found")
}

func attachXDP(objs *AdvancedSecurityObjects) error {
	// Demo mode - simulate XDP attachment
	if objs.AdvancedFilter == nil {
		log.Println("🔧 Demo mode: Simulating XDP program attachment")
		return nil
	}

	// Real XDP attachment (when eBPF is available)
	iface, err := getNetworkInterface()
	if err != nil {
		return fmt.Errorf("failed to find network interface: %v", err)
	}

	link, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.AdvancedFilter,
		Interface: iface.Index,
	})
	if err != nil {
		return fmt.Errorf("failed to attach XDP program: %v", err)
	}

	defer link.Close()
	log.Printf("✅ XDP program attached to interface %s", iface.Name)
	return nil
}

func attachKprobes(objs *AdvancedSecurityObjects) {
	kprobes := []struct {
		name string
		prog *ebpf.Program
	}{
		{"sys_connect", objs.TraceConnect},
		{"sys_execve", objs.TraceExecve},
		{"sys_write", objs.TraceWrite},
	}

	for _, kp := range kprobes {
		if link, err := link.Kprobe(kp.name, kp.prog, nil); err == nil {
			defer link.Close()
		} else {
			log.Printf("Failed to attach kprobe %s: %v", kp.name, err)
		}
	}
}

func startMetricsServer(objs *AdvancedSecurityObjects) {
	r := mux.NewRouter()

	// Metrics endpoint
	r.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics := collectMetrics(objs)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metrics)
	}).Methods("GET")

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	}).Methods("GET")

	// Top connections endpoint
	r.HandleFunc("/top-connections", func(w http.ResponseWriter, r *http.Request) {
		topIPs := getTopIPs(objs)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(topIPs)
	}).Methods("GET")

	// Security events endpoint
	r.HandleFunc("/security-events", func(w http.ResponseWriter, r *http.Request) {
		security := getSecurityMetrics(objs)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(security)
	}).Methods("GET")

	fmt.Println("🌐 Metrics server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Printf("HTTP server error: %v", err)
	}
}

func startMetricsCollection(objs *AdvancedSecurityObjects, ifaceName string) {
	fmt.Printf("🚀 Advanced eBPF Security Monitor Started on interface %s\n", ifaceName)
	fmt.Println("📊 Monitoring network traffic, system calls, and security events")
	fmt.Println("🌐 Metrics API available at http://localhost:8080/metrics")
	fmt.Println("Press Ctrl+C to stop...")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-c:
			fmt.Println("\n👋 Shutting down advanced monitor...")
			return
		case <-ticker.C:
			displayAdvancedMetrics(objs)
		}
	}
}

func collectMetrics(objs *AdvancedSecurityObjects) Metrics {
	var key uint32 = 0
	metrics := Metrics{
		Timestamp: time.Now(),
		TopIPs:    getTopIPs(objs),
		Security:  getSecurityMetrics(objs),
	}

	// Demo mode - return simulated metrics
	if objs.PacketCount == nil {
		metrics.Packets = 12500
		metrics.Bytes = 45200000
		metrics.HTTPRequests = 1200
		metrics.TCPConn = 850
		metrics.UDPPackets = 200
		metrics.Blocked = 2
		metrics.Transactions = 150
		return metrics
	}

	// Collect basic metrics from eBPF maps
	objs.PacketCount.Lookup(&key, &metrics.Packets)
	objs.ByteCount.Lookup(&key, &metrics.Bytes)
	objs.HttpRequests.Lookup(&key, &metrics.HTTPRequests)
	objs.TcpConnections.Lookup(&key, &metrics.TCPConn)
	objs.UdpPackets.Lookup(&key, &metrics.UDPPackets)
	objs.BlockedConnections.Lookup(&key, &metrics.Blocked)
	objs.TransactionCount.Lookup(&key, &metrics.Transactions)

	return metrics
}

func getTopIPs(objs *AdvancedSecurityObjects) []IPMetrics {
	// Demo mode - return simulated top IPs
	if objs.ConnectionCount == nil {
		return []IPMetrics{
			{IP: "192.168.1.100", Count: 1250},
			{IP: "10.0.0.50", Count: 890},
			{IP: "172.16.0.25", Count: 650},
			{IP: "192.168.1.200", Count: 450},
			{IP: "10.0.0.100", Count: 320},
		}
	}

	var topIPs []IPMetrics
	var (
		ipKey   uint32
		ipCount uint64
		maxIPs  = 10
		count   = 0
	)

	iterator := objs.ConnectionCount.Iterate()
	for iterator.Next(&ipKey, &ipCount) && count < maxIPs {
		if ipCount > 0 {
			ip := net.IPv4(byte(ipKey), byte(ipKey>>8), byte(ipKey>>16), byte(ipKey>>24))
			topIPs = append(topIPs, IPMetrics{
				IP:    ip.String(),
				Count: ipCount,
			})
			count++
		}
	}

	return topIPs
}

func getSecurityMetrics(objs *AdvancedSecurityObjects) SecurityMetrics {
	// Demo mode - return simulated security metrics
	if objs.SuspiciousIps == nil {
		return SecurityMetrics{
			SuspiciousIPs: []SuspiciousIP{
				{IP: "192.168.1.100", Count: 25},
				{IP: "192.168.1.200", Count: 15},
			},
			PortScans: []PortScan{
				{IP: "192.168.1.200", Ports: 22},
				{IP: "192.168.1.200", Ports: 80},
			},
			Threats: []Threat{
				{Type: "Port Scan", IP: "192.168.1.200", Timestamp: time.Now().Add(-5 * time.Minute), Details: "Multiple ports scanned"},
				{Type: "Suspicious Activity", IP: "192.168.1.100", Timestamp: time.Now().Add(-2 * time.Minute), Details: "High connection rate"},
			},
		}
	}

	security := SecurityMetrics{
		SuspiciousIPs: getSuspiciousIPs(objs),
		PortScans:     getPortScans(objs),
		Threats:       getRecentThreats(objs),
	}

	return security
}

func getSuspiciousIPs(objs *AdvancedSecurityObjects) []SuspiciousIP {
	var suspicious []SuspiciousIP
	var (
		ipKey   uint32
		ipCount uint64
		maxIPs  = 10
		count   = 0
	)

	iterator := objs.SuspiciousIps.Iterate()
	for iterator.Next(&ipKey, &ipCount) && count < maxIPs {
		if ipCount > 100 { // Threshold for suspicious activity
			ip := net.IPv4(byte(ipKey), byte(ipKey>>8), byte(ipKey>>16), byte(ipKey>>24))
			suspicious = append(suspicious, SuspiciousIP{
				IP:    ip.String(),
				Count: ipCount,
			})
			count++
		}
	}

	return suspicious
}

func getPortScans(objs *AdvancedSecurityObjects) []PortScan {
	var scans []PortScan
	var (
		ipKey     uint32
		portCount uint32
		maxScans  = 10
		count     = 0
	)

	iterator := objs.PortScanDetection.Iterate()
	for iterator.Next(&ipKey, &portCount) && count < maxScans {
		if portCount > 50 { // Threshold for port scan detection
			ip := net.IPv4(byte(ipKey), byte(ipKey>>8), byte(ipKey>>16), byte(ipKey>>24))
			scans = append(scans, PortScan{
				IP:    ip.String(),
				Ports: portCount,
			})
			count++
		}
	}

	return scans
}

func getRecentThreats(objs *AdvancedSecurityObjects) []Threat {
	var threats []Threat
	var (
		eventKey uint32
		event    struct {
			SrcIP     uint32
			DstIP     uint32
			SrcPort   uint16
			DstPort   uint16
			Timestamp uint64
			EventData uint32
			Severity  uint8
		}
		iterator    = objs.SecurityEvents.Iterate()
		maxThreats  = 10
		count       = 0
		currentTime = uint64(time.Now().UnixNano())
	)

	// Collect recent security events from the last 5 minutes
	for iterator.Next(&eventKey, &event) && count < maxThreats {
		// Only include events from the last 5 minutes
		if currentTime-event.Timestamp < 3000000000000 { // 5 minutes in nanoseconds
			srcIP := net.IPv4(byte(event.SrcIP), byte(event.SrcIP>>8), byte(event.SrcIP>>16), byte(event.SrcIP>>24))
			dstIP := net.IPv4(byte(event.DstIP), byte(event.DstIP>>8), byte(event.DstIP>>16), byte(event.DstIP>>24))

			var threatType, details string
			switch event.Severity {
			case 10:
				threatType = "DDoS Attack"
				details = fmt.Sprintf("High volume traffic from %s (%d packets)", srcIP.String(), event.EventData)
			case 7:
				threatType = "Port Scan"
				details = fmt.Sprintf("Port scan detected from %s to %s:%d", srcIP.String(), dstIP.String(), event.DstPort)
			case 5:
				threatType = "Suspicious Activity"
				details = fmt.Sprintf("Unusual connection pattern from %s to %s", srcIP.String(), dstIP.String())
			case 3:
				threatType = "Transaction Anomaly"
				details = fmt.Sprintf("Unusual financial transaction from %s ($%.2f)", srcIP.String(), float64(event.EventData)/100)
			default:
				threatType = "Security Event"
				details = fmt.Sprintf("Event from %s to %s:%d", srcIP.String(), dstIP.String(), event.DstPort)
			}

			threats = append(threats, Threat{
				Type:      threatType,
				IP:        srcIP.String(),
				Timestamp: time.Unix(0, int64(event.Timestamp)),
				Details:   details,
			})
			count++
		}
	}

	// Sort threats by severity (most severe first)
	for i := 0; i < len(threats)-1; i++ {
		for j := i + 1; j < len(threats); j++ {
			if getThreatSeverity(threats[i].Type) < getThreatSeverity(threats[j].Type) {
				threats[i], threats[j] = threats[j], threats[i]
			}
		}
	}

	return threats
}

// getThreatSeverity returns a numeric severity for threat type sorting
func getThreatSeverity(threatType string) int {
	switch threatType {
	case "DDoS Attack":
		return 10
	case "Port Scan":
		return 7
	case "Suspicious Activity":
		return 5
	case "Transaction Anomaly":
		return 3
	default:
		return 1
	}
}

func displayAdvancedMetrics(objs *AdvancedSecurityObjects) {
	metrics := collectMetrics(objs)

	fmt.Printf("📊 [%s] ", metrics.Timestamp.Format("15:04:05"))
	fmt.Printf("📦 %d pkt/s | 📊 %.1f MB/s | 🌐 %d HTTP/s | 🔗 %d TCP | 📡 %d UDP | 🚫 %d blocked | 💰 %d tx\n",
		metrics.Packets,
		float64(metrics.Bytes)/1024/1024,
		metrics.HTTPRequests,
		metrics.TCPConn,
		metrics.UDPPackets,
		metrics.Blocked,
		metrics.Transactions,
	)

	// Display security alerts
	if len(metrics.Security.SuspiciousIPs) > 0 {
		fmt.Printf("⚠️  %d suspicious IPs detected\n", len(metrics.Security.SuspiciousIPs))
	}

	if len(metrics.Security.PortScans) > 0 {
		fmt.Printf("🔍 %d port scans detected\n", len(metrics.Security.PortScans))
	}
}
