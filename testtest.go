package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"net"
	"syscall"

	"golang.org/x/net/ipv4"
	"golang.org/x/time/rate"
)

const (
	UDP_SEGMENT = 103
)

func getCmsg(size int) []byte {
	b := make([]byte, 8+4+4+2)
	binary.LittleEndian.PutUint64(b[:8], 18)
	binary.LittleEndian.PutUint32(b[8:12], uint32(syscall.IPPROTO_UDP))
	binary.LittleEndian.PutUint32(b[12:16], uint32(UDP_SEGMENT))
	binary.LittleEndian.PutUint16(b[16:], uint16(size))
	return b
}

func udpReadBatch() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(192, 168, 1, 23),
		Port: 22223,
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	packetconn := ipv4.NewPacketConn(conn)
	// Fd, _ := conn.File()
	// fd := int(Fd.Fd())
	// if err := syscall.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_ZEROCOPY, 1); err != nil {
	// 	log.Fatalf("Failed to set SO_ZEROCOPY: %v", err)
	// }

	msgs := make([]ipv4.Message, 128)
	for v := range msgs {
		msgs[v].Buffers = make([][]byte, 1)
		for y := range msgs[v].Buffers {
			msgs[v].Buffers[y] = make([]byte, 1500)
		} //
	}
	count := 0
	go func() {
		for {
			time.Sleep(time.Second * 3)
			fmt.Println("count---------", count)
			count = 0
		}

	}()
	for {
		n, err := packetconn.ReadBatch(msgs, 0)
		if err != nil {
			fmt.Println(err.Error())
		}

		count += n

	}
}

func main() {

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 1, 100),
		Port: 12121,
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	packetconn := ipv4.NewPacketConn(conn)
	size := 18
	b := make([]byte, size*35)
	msgs := make([]ipv4.Message, 1)
	for v := range msgs {
		msgs[v].Buffers = make([][]byte, 1)
		for y := range msgs[v].Buffers {
			msgs[v].Buffers[y] = b
		}
		msgs[v].OOB = getCmsg(size)
		//
	}
	ctx := context.Background()
	r := rate.NewLimiter(1024*1024*25, 1024*1024*300)
	for i := 0; i < 3000000; i++ {
		// conn.Write(b)
		n, err := packetconn.WriteBatch(msgs, 0)
		r.WaitN(ctx, n*size*35)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
	}

}
