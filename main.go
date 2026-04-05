package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
	"github.com/urfave/cli/v2"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang Monitor ./bpf/monitor.c
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang Security ./bpf/security.c

func main() {
	app := &cli.App{
		Name:  "ebpf-demo",
		Usage: "eBPF monitoring and security demo",
		Commands: []*cli.Command{
			{
				Name:  "monitor",
				Usage: "Start network monitoring",
				Action: func(c *cli.Context) error {
					return runMonitor()
				},
			},
			{
				Name:  "security",
				Usage: "Start security monitoring",
				Action: func(c *cli.Context) error {
					return runSecurity()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runMonitor() error {
	// Remove memory limit for eBPF
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatalf("Failed to remove memory limit: %v", err)
	}

	// Load pre-compiled eBPF program
	objs := MonitorObjects{}
	if err := LoadMonitorObjects(&objs, nil); err != nil {
		log.Fatalf("Loading objects failed: %v", err)
	}
	defer objs.Close()

	// Get the first network interface
	iface, err := net.InterfaceByName("eth0")
	if err != nil {
		// Try other common interface names
		iface, err = net.InterfaceByName("en0")
		if err != nil {
			iface, err = net.InterfaceByName("lo")
			if err != nil {
				log.Fatalf("Failed to find network interface: %v", err)
			}
		}
	}

	// Attach XDP program
	xdpLink, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.MonitorNetwork,
		Interface: iface.Index,
	})
	if err != nil {
		log.Fatalf("Failed to attach XDP program: %v", err)
	}
	defer xdpLink.Close()

	// Attach kprobes
	kprobeConnect, err := link.Kprobe("sys_connect", objs.TraceConnect, nil)
	if err != nil {
		log.Printf("Failed to attach kprobe sys_connect: %v", err)
	} else {
		defer kprobeConnect.Close()
	}

	kprobeExecve, err := link.Kprobe("sys_execve", objs.TraceExecve, nil)
	if err != nil {
		log.Printf("Failed to attach kprobe sys_execve: %v", err)
	} else {
		defer kprobeExecve.Close()
	}

	tracepointWrite, err := link.Tracepoint("syscalls", "sys_enter_write", objs.TraceWrite, nil)
	if err != nil {
		log.Printf("Failed to attach tracepoint sys_enter_write: %v", err)
	} else {
		defer tracepointWrite.Close()
	}

	fmt.Printf("🚀 eBPF Network Monitor Started on interface %s\n", iface.Name)
	fmt.Println("📊 Monitoring network traffic, system calls, and HTTP requests")
	fmt.Println("Press Ctrl+C to stop...")

	// Start metrics collection
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-c:
			fmt.Println("\n👋 Shutting down monitor...")
			return nil
		case <-ticker.C:
			displayMetrics(&objs)
		}
	}
}

func runSecurity() error {
	// Remove memory limit for eBPF
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatalf("Failed to remove memory limit: %v", err)
	}

	// Load pre-compiled eBPF program
	objs := SecurityObjects{}
	if err := LoadSecurityObjects(&objs, nil); err != nil {
		log.Fatalf("Loading objects failed: %v", err)
	}
	defer objs.Close()

	// Get the first network interface
	iface, err := net.InterfaceByName("eth0")
	if err != nil {
		// Try other common interface names
		iface, err = net.InterfaceByName("en0")
		if err != nil {
			iface, err = net.InterfaceByName("lo")
			if err != nil {
				log.Fatalf("Failed to find network interface: %v", err)
			}
		}
	}

	// Attach XDP program for security filtering
	xdpLink, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.FilterTransactions,
		Interface: iface.Index,
	})
	if err != nil {
		log.Fatalf("Failed to attach XDP program: %v", err)
	}
	defer xdpLink.Close()

	// Attach security kprobes
	kprobeExecve, err := link.Kprobe("sys_execve", objs.SecurityMonitor, nil)
	if err != nil {
		log.Printf("Failed to attach kprobe sys_execve: %v", err)
	} else {
		defer kprobeExecve.Close()
	}

	kprobeConnect, err := link.Kprobe("sys_connect", objs.MonitorConnections, nil)
	if err != nil {
		log.Printf("Failed to attach kprobe sys_connect: %v", err)
	} else {
		defer kprobeConnect.Close()
	}

	fmt.Printf("🔒 eBPF Security Monitor Started on interface %s\n", iface.Name)
	fmt.Println("🛡️  Monitoring for suspicious activity and blocking threats")
	fmt.Println("💰 Tracking financial transactions")
	fmt.Println("Press Ctrl+C to stop...")

	// Start security metrics collection
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-c:
			fmt.Println("\n👋 Shutting down security monitor...")
			return nil
		case <-ticker.C:
			displaySecurityMetrics(&objs)
		}
	}
}

func displayMetrics(objs *MonitorObjects) {
	var key uint32 = 0

	// Get packet count
	var packetCount uint64
	if err := objs.PacketCount.Lookup(&key, &packetCount); err == nil {
		fmt.Printf("📦 Packets/sec: %d", packetCount)
	}

	// Get byte count
	var byteCount uint64
	if err := objs.ByteCount.Lookup(&key, &byteCount); err == nil {
		fmt.Printf(" | 📊 MB/sec: %.2f", float64(byteCount)/1024/1024)
	}

	// Get HTTP requests
	var httpRequests uint64
	if err := objs.HttpRequests.Lookup(&key, &httpRequests); err == nil {
		fmt.Printf(" | 🌐 HTTP/sec: %d", httpRequests)
	}

	// Get TCP connections
	var tcpConnections uint64
	if err := objs.TcpConnections.Lookup(&key, &tcpConnections); err == nil {
		fmt.Printf(" | 🔗 TCP: %d", tcpConnections)
	}

	// Get UDP packets
	var udpPackets uint64
	if err := objs.UdpPackets.Lookup(&key, &udpPackets); err == nil {
		fmt.Printf(" | 📡 UDP: %d", udpPackets)
	}

	fmt.Println()

	// Display top connections
	fmt.Printf("🔍 Top connections:\n")
	var (
		connKey   uint32
		connCount uint64
		maxConn   uint64
		topIP     uint32
		iterator  = objs.ConnectionCount.Iterate()
	)

	for iterator.Next(&connKey, &connCount) {
		if connCount > maxConn {
			maxConn = connCount
			topIP = connKey
		}
	}

	if topIP != 0 {
		ip := net.IPv4(byte(topIP), byte(topIP>>8), byte(topIP>>16), byte(topIP>>24))
		fmt.Printf("   📍 %s: %d connections\n", ip.String(), maxConn)
	}
}

func displaySecurityMetrics(objs *SecurityObjects) {
	var key uint32 = 0

	// Get transaction count
	var transactionCount uint64
	if err := objs.TransactionCount.Lookup(&key, &transactionCount); err == nil {
		fmt.Printf("💰 Transactions/sec: %d", transactionCount)
	}

	// Get blocked connections
	var blockedConnections uint64
	if err := objs.BlockedConnections.Lookup(&key, &blockedConnections); err == nil {
		fmt.Printf(" | 🚫 Blocked: %d", blockedConnections)
	}

	fmt.Println()

	// Display suspicious IPs
	fmt.Printf("🔍 Security events:\n")
	var (
		suspKey   uint32
		suspCount uint64
		iterator  = objs.SuspiciousIps.Iterate()
	)

	for iterator.Next(&suspKey, &suspCount) {
		if suspCount > 10 {
			ip := net.IPv4(byte(suspKey), byte(suspKey>>8), byte(suspKey>>16), byte(suspKey>>24))
			fmt.Printf("   ⚠️  Suspicious IP %s: %d events\n", ip.String(), suspCount)
		}
	}

	// Display recent transactions
	fmt.Printf("💸 Recent transactions:\n")
	var (
		eventKey uint32
		event    struct {
			SrcIP     uint32
			DstIP     uint32
			SrcPort   uint16
			DstPort   uint16
			Timestamp uint64
			Amount    uint32
			Protocol  uint8
		}
		eventIterator = objs.TransactionEvents.Iterate()
		count         int
	)

	for eventIterator.Next(&eventKey, &event) && count < 5 {
		srcIP := net.IPv4(byte(event.SrcIP), byte(event.SrcIP>>8), byte(event.SrcIP>>16), byte(event.SrcIP>>24))
		dstIP := net.IPv4(byte(event.DstIP), byte(event.DstIP>>8), byte(event.DstIP>>16), byte(event.DstIP>>24))
		fmt.Printf("   💰 $%.2f %s:%d → %s:%d\n",
			float64(event.Amount)/100, srcIP.String(), event.SrcPort,
			dstIP.String(), event.DstPort)
		count++
	}
}
