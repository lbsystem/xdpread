#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/udp.h>

struct {
    __uint(type, BPF_MAP_TYPE_XSKMAP);
    __uint(key_size, sizeof(__u32));
    __uint(value_size, sizeof(__u32));
    __uint(max_entries, 64);
} xsks_map SEC(".maps");

static volatile unsigned const short PORT;

struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(__u32));
    __uint(value_size, sizeof(__u32));
    __uint(max_entries, 64);
} events SEC(".maps");

SEC("xdp")
int xdp_filter_port(struct xdp_md *ctx) {
    void *data_end = (void *)(long)ctx->data_end;
    void *data = (void *)(long)ctx->data;
    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end)
        return XDP_PASS;

    // 确保包是 IP 包
    if (eth->h_proto != bpf_htons(ETH_P_IP))
        return XDP_PASS;

    struct iphdr *ip = (struct iphdr *)(eth + 1);
    if ((void *)(ip + 1) > data_end || ip->protocol != 17)
        return XDP_PASS;

    struct udphdr *udp = (struct udphdr *)(ip + 1);
    if ((void *)(udp + 1) > data_end)
        return XDP_PASS;

    // 检查 UDP 目的端口
    if (udp->dest == bpf_htons(PORT)) {
        int index = ctx->rx_queue_index;
        bpf_printk("Hello, world, from BPF! My PID is %d %d\n", 12346, index);
        return bpf_redirect_map(&xsks_map, index, 0);
    }

    return XDP_PASS;
}

char _license[] SEC("license") = "GPL";
