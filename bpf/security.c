#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <linux/udp.h>
#include <linux/in.h>
#include <bpf/bpf_endian.h>

struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_HASH);
    __type(key, __u32);
    __type(value, __u64);
    __uint(max_entries, 1024);
} suspicious_ips SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __type(key, __u32);
    __type(value, __u64);
    __uint(max_entries, 256);
} blocked_connections SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __type(key, __u32);
    __type(value, __u64);
    __uint(max_entries, 256);
} transaction_count SEC(".maps");

struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_HASH);
    __type(key, __u32);
    __type(value, __u64);
    __uint(max_entries, 1024);
} port_scan_detection SEC(".maps");

struct transaction_event {
    __u32 src_ip;
    __u32 dst_ip;
    __u16 src_port;
    __u16 dst_port;
    __u64 timestamp;
    __u32 amount;  // Simulated amount in cents
    __u8 protocol;
};

struct {
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __type(key, __u32);
    __type(value, struct transaction_event);
    __uint(max_entries, 10000);
} transaction_events SEC(".maps");

SEC("xdp")
int filter_transactions(struct xdp_md *ctx) {
    void *data_end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;
    
    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end)
        return XDP_PASS;
    
    if (eth->h_proto != __constant_htons(ETH_P_IP))
        return XDP_PASS;
    
    struct iphdr *ip = (void *)(eth + 1);
    if ((void *)(ip + 1) > data_end)
        return XDP_PASS;
    
    __u32 src_ip = ip->saddr;
    __u32 dst_ip = ip->daddr;
    
    // Check if source IP is suspicious
    __u64 *suspicious_count = bpf_map_lookup_elem(&suspicious_ips, &src_ip);
    if (suspicious_count && *suspicious_count > 100) {
        // Block suspicious traffic
        __u32 key = 0;
        __u64 *blocked = bpf_map_lookup_elem(&blocked_connections, &key);
        if (blocked) {
            __sync_fetch_and_add(blocked, 1);
        }
        return XDP_DROP;
    }
    
    // Track port scanning
    if (ip->protocol == IPPROTO_TCP) {
        struct tcphdr *tcp = (void *)(ip + 1);
        if ((void *)(tcp + 1) > data_end)
            return XDP_PASS;
        
        __u16 dst_port = bpf_ntohs(tcp->dest);
        
        // Track unique destination ports per source IP
        __u64 *port_count = bpf_map_lookup_elem(&port_scan_detection, &src_ip);
        if (port_count) {
            __sync_fetch_and_add(port_count, 1);
            if (*port_count > 50) {
                // Potential port scan - block
                __u32 key = 0;
                __u64 *blocked = bpf_map_lookup_elem(&blocked_connections, &key);
                if (blocked) {
                    __sync_fetch_and_add(blocked, 1);
                }
                return XDP_DROP;
            }
        } else {
            __u64 new_count = 1;
            bpf_map_update_elem(&port_scan_detection, &src_ip, &new_count, BPF_NOEXIST);
        }
        
        // Monitor financial transactions (ports 443, 8443)
        if (dst_port == 443 || dst_port == 8443) {
            __u32 key = 0;
            __u64 *count = bpf_map_lookup_elem(&transaction_count, &key);
            if (count) {
                __sync_fetch_and_add(count, 1);
            }
            
            // Store transaction event
            struct transaction_event event = {
                .src_ip = src_ip,
                .dst_ip = dst_ip,
                .src_port = bpf_ntohs(tcp->source),
                .dst_port = dst_port,
                .timestamp = bpf_ktime_get_ns(),
                .amount = bpf_get_prandom_u32() % 100000,  // Random amount for demo
                .protocol = IPPROTO_TCP
            };
            
            __u32 event_key = (__u32)(bpf_ktime_get_ns() % 10000);
            bpf_map_update_elem(&transaction_events, &event_key, &event, BPF_ANY);
        }
    }
    
    return XDP_PASS;
}

// Security monitoring for process execution
SEC("kprobe/sys_execve")
int security_monitor(struct pt_regs *ctx) {
    // Track suspicious process executions
    __u32 key = 0;
    __u64 *count = bpf_map_lookup_elem(&blocked_connections, &key);
    if (count) {
        __sync_fetch_and_add(count, 1);
    }
    
    return 0;
}

// Monitor network connections for security
SEC("kprobe/sys_connect")
int monitor_connections(struct pt_regs *ctx) {
    __u32 key = 0;
    __u64 *count = bpf_map_lookup_elem(&transaction_count, &key);
    if (count) {
        __sync_fetch_and_add(count, 1);
    }
    
    return 0;
}

char _license[] SEC("license") = "GPL";
