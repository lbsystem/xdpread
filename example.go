package main

import (
	"fmt"
	"log"
	// "os"
	// "os/signal"
	// "runtime/pprof"
	// "syscall"

	"time"

	"sync/atomic"

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

	xsk, bpf, l := udpbpf.NewAFXDP("enp6s16", 12121, 0, true)
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
			c := atomic.LoadInt64(&count)
			if c == 0 {
				fmt.Printf("\r                            ")
				continue
			}
			atomic.SwapInt64(&count, 0)
			fmt.Printf("\rspeed is %dM", c/1024/1024)
		}
	}()
	ff := xsk.HandleRecv()
	for {
		n := ff(b[:0])
		for i := 0; i < n; i++ {
			_, data, err := udpbpf.DecodeUdp(b[i])
			atomic.AddInt64(&count, int64(len(data)))
			if err != nil {
				continue
			}

		}
	}
}
