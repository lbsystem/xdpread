package main

import (
	"fmt"
	

	udpbpf "github.com/lbsystem/xdpread/ebpf"
)

func main() {
	
	
	xsk, bpf, l := udpbpf.NewAFXDP("enp6s16", 12121, 0)

	defer bpf.Close()
	defer l.Close()
	b := make([][]byte, 100)
	for {
	
		n := xsk.Read(b)
		for i := 0; i < n; i++ {

			src, dst, data, err := udpbpf.DecodeUdp(b[i])
			fmt.Println(src, dst, data, err)
			

		}

	}
}
