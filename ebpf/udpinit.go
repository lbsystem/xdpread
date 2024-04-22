package udpbpf

import (
	"C"
)
import (
	"encoding/binary"
	"fmt"
	"log"
	"net"

	"github.com/cilium/ebpf/link"
	"github.com/lixiangzhong/xdp"
	"golang.org/x/sys/unix"
)

func newEbpf(ifname string, port uint16) (*udpObjects, link.Link, *net.Interface) {
	var bpf1 udpObjects
	collectionSpec, err := loadUdp()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = collectionSpec.RewriteConstants(map[string]interface{}{"PORT": port})
	if err != nil {
		log.Fatal(err.Error())
	}
	err = collectionSpec.LoadAndAssign(&bpf1, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	iface1, err := net.InterfaceByName(ifname)
	if err != nil {
		log.Fatal(err)
	}

	l, err := link.AttachXDP(link.XDPOptions{
		Program:   bpf1.XdpFilterPort,
		Interface: iface1.Index,
		Flags:     unix.XDP_FLAGS_DRV_MODE,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	return &bpf1, l, iface1

}

func createsocket(ifidx int, queueid int, poll bool) *xdp.Socket {
	umem, err := xdp.NewUmem(&xdp.UmemConfig{
		FillSize: xdp.DEFAULT_FILL_SIZE,
		CompSize: xdp.DEFAULT_COMP_SIZE,
		Size:     uint32(16 << 20),
	})
	if err != nil {
		log.Fatal(err)
	}

	xsk, err := xdp.NewSocket(ifidx, umem, &xdp.SocketConfig{
		RxSize:    xdp.DEFAULT_RX_SIZE,
		TxSize:    xdp.DEFAULT_TX_SIZE,
		QueueID:   queueid,
		BindFlags: unix.XDP_ZEROCOPY, //XDP_ZEROCOPY XDP_USE_NEED_WAKEUP XDP_SHARED_UMEM XDP_COPY
		Poll:      poll,
	})
	if err != nil {
		log.Fatal(err)
	}

	return xsk

}

func NewAFXDP(ifname string, port, queueID int, poll bool) (*xdp.Socket, *udpObjects, link.Link) {

	bpf, l, link := newEbpf(ifname, uint16(port))

	xsk := createsocket(link.Index, queueID, poll)

	err := bpf.XsksMap.Put(uint32(queueID), uint32(xsk.FD()))
	if err != nil {
		log.Fatal(err.Error())
	}

	return xsk, bpf, l
}
func DecodeUdp(b []byte) (src *net.UDPAddr, data []byte, e error) {

	if len(b) < 42 {
		return nil, nil, fmt.Errorf("It is not an IP data")
	}
	// IPlength := binary.BigEndian.Uint16(b[16:18])
	IPlength := uint16(b[17]) | uint16(b[16])<<8
	// IPlength:=28+18
	data1 := b[14+28 : 14+IPlength]
	srcIP := net.UDPAddr{
		IP: net.IP(b[26:30]),
		// Port: *(*int)(unsafe.Pointer(&b[20])),
		Port: int(binary.BigEndian.Uint16(b[14+20:])),
	}
	// dstIP := net.UDPAddr{
	// 	IP:   net.IP(b[30:34]),
	// 	Port: int(binary.BigEndian.Uint16(b[14+22:])),
	// }
	return &srcIP, data1, nil
}
