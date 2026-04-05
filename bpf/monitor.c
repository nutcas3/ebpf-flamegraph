#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <linux/udp.h>
#include <linux/in.h>
#include <bpf/bpf_endian.h>

// Maps for storing metrics
struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __type(key, __u32);
    __type(value, __u64);
    __uint(max_entries, 256);
} packet_count SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __type(key, __u32);
    __type(value, __u64);
    __uint(max_entries, 256);
} byte_count SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_HASH);
    __type(key, __u32);
    __type(value, __u64);
    __uint(max_entries, 1024);
} connection_count SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __type(key, __u32);
    __type(value, __u64);
    __uint(max_entries, 256);
} http_requests SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __type(key, __u32);
    __type(value, __u64);
    __uint(max_entries, 256);
} tcp_connections SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __type(key, __u32);
    __type(value, __u64);
    __uint(max_entries, 256);
} udp_packets SEC(".maps");

// Packet processing function
SEC("xdp")
int monitor_network(struct xdp_md *ctx) {
    void *data_end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;
    
    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end)
        return XDP_PASS;
    
    // Only process IP packets
    if (eth->h_proto != __constant_htons(ETH_P_IP))
        return XDP_PASS;
    
    struct iphdr *ip = (void *)(eth + 1);
    if ((void *)(ip + 1) > data_end)
        return XDP_PASS;
    
    __u32 key = 0;
    __u64 *count;
    __u64 packet_len = bpf_ntohs(ip->tot_len);
    
    // Update packet count
    count = bpf_map_lookup_elem(&packet_count, &key);
    if (count) {
        __sync_fetch_and_add(count, 1);
    }
    
    // Update byte count
    count = bpf_map_lookup_elem(&byte_count, &key);
    if (count) {
        __sync_fetch_and_add(count, packet_len);
    }
    
    // Track connections by source IP
    __u32 src_ip = ip->saddr;
    count = bpf_map_lookup_elem(&connection_count, &src_ip);
    if (count) {
        __sync_fetch_and_add(count, 1);
    } else {
        __u64 new_count = 1;
        bpf_map_update_elem(&connection_count, &src_ip, &new_count, BPF_NOEXIST);
    }
    
    // Protocol-specific tracking
    if (ip->protocol == IPPROTO_TCP) {
        struct tcphdr *tcp = (void *)(ip + 1);
        if ((void *)(tcp + 1) > data_end)
            return XDP_PASS;
        
        // Track TCP connections
        key = 0;
        count = bpf_map_lookup_elem(&tcp_connections, &key);
        if (count) {
            __sync_fetch_and_add(count, 1);
        }
        
        // Check for HTTP traffic (port 80 or 8080)
        __u16 dport = bpf_ntohs(tcp->dest);
        if (dport == 80 || dport == 8080) {
            key = 0;
            count = bpf_map_lookup_elem(&http_requests, &key);
            if (count) {
                __sync_fetch_and_add(count, 1);
            }
        }
    } else if (ip->protocol == IPPROTO_UDP) {
        // Track UDP packets
        key = 0;
        count = bpf_map_lookup_elem(&udp_packets, &key);
        if (count) {
            __sync_fetch_and_add(count, 1);
        }
    }
    
    return XDP_PASS;
}

// System call monitoring for security
SEC("kprobe/sys_connect")
int trace_connect(struct pt_regs *ctx) {
    __u32 key = 0;
    __u64 *count;
    
    // Track connection attempts
    count = bpf_map_lookup_elem(&tcp_connections, &key);
    if (count) {
        __sync_fetch_and_add(count, 1);
    }
    
    return 0;
}

SEC("kprobe/sys_execve")
int trace_execve(struct pt_regs *ctx) {
    __u32 key = 0;
    __u64 *count;
    
    // Track process executions
    count = bpf_map_lookup_elem(&connection_count, &key);
    if (count) {
        __sync_fetch_and_add(count, 1);
    }
    
    return 0;
}

SEC("tracepoint/syscalls/sys_enter_write")
int trace_write(struct pt_regs *ctx) {
    __u32 key = 0;
    __u64 *count;
    
    // Track write operations (potential HTTP responses)
    count = bpf_map_lookup_elem(&http_requests, &key);
    if (count) {
        __sync_fetch_and_add(count, 1);
    }
    
    return 0;
}

char _license[] SEC("license") = "GPL";
