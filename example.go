package main

import (
	"fmt"
	"log"

	udpbpf "github.com/lbsystem/xdpread/ebpf"
	"github.com/lixiangzhong/xdp"
)

func main() {
	
	
	xsk, bpf, l := udpbpf.NewAFXDP("enp6s16", 12121, 0)
	n, _, err := xdp.GetNicQueues("enp6s16")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(n)

	defer bpf.Close()
	defer l.Close()
	// select{}
	b := make([][]byte,2048)
	// rmsgs := make([]ipv4.Message, 128)

	for {	
		n:=xsk.HandleRecv(b[:0])
	
		for i := 0; i < n; i++ {
			_, _, _, err := udpbpf.DecodeUdp(b[i])		
			if err!=nil{
				continue
			}
			// fmt.Println(a.String())
		}
	}
}
