package main

import (
	"fmt"
	"log"

	"time"

	"sync/atomic"

	"github.com/lbsystem/protocol"
	udpbpf "github.com/lbsystem/xdpread/ebpf"
	"github.com/lixiangzhong/xdp"
)

func main() {
	// f, err := os.Create("/goProject/xdp/default.prof")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// time.AfterFunc(time.Second*2, func() {
	// 	if err := pprof.StartCPUProfile(f); err != nil {
	// 		panic(err)
	// 	}
	// })

	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	// go func() {
	// 	<-signalChan

	// 	pprof.StopCPUProfile()
	// 	f.Close()
	// 	fmt.Println("file close")
	// 	os.Exit(0)
	// }()

	xsk, bpf, l := udpbpf.NewAFXDP("enp6s16", 12121, 0, false)
	n, _, err := xdp.GetNicQueues("enp6s16")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(n)

	defer bpf.Close()
	defer l.Close()
	// select{}
	b := make([][]byte, 2048)
	// rmsgs := make([]ipv4.Message, 128)
	var count int64
	go func() {
		for {
			time.Sleep(time.Second)
			c := atomic.SwapInt64(&count, 0)
			if c == 0 {
				fmt.Printf("\r                            ")
				continue
			}

			fmt.Printf("\rspeed is %dM", c/1024/1024)
		}
	}()
	var openInterrupt bool = true
	xdpRead := xsk.HandleRecv(&openInterrupt)
	// var waitTime = 1
	eth := protocol.NewEthernet()
	for {
		n, s := xdpRead(b[:0])
		// if n == 0 {
		// 	// time.Sleep(time.Microsecond * 5 * time.Duration(waitTime))
		// 	if waitTime < 128 {
		// 		waitTime++
		// 	} else {
		// 		openInterrupt = true
		// 	}
		// 	continue
		// }
		// openInterrupt = false
		// waitTime = 1
		for i := 0; i < n; i++ {
			// _, data, err := udpbpf.DecodeUdp(b[i])
			err := eth.UnmarshalBinary(b[i])
			if err != nil {
				fmt.Println(err.Error())
			}
			ipdata, ok := eth.Data.(*protocol.IPv4)
			if ok {
				udpdata, ok := ipdata.Data.(*protocol.UDP)
				if ok {
					atomic.AddInt64(&count, int64(udpdata.Len()-8))

				}
			}
		}
		if n > 0 {
			s.Submit_cons(uint32(n))
		}

	}
}
