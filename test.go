package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"sync/atomic"
)

func main() {
	// f, err := os.Create("/goProject/test/xdpread/cpu.prof")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
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

	u, err := net.ResolveUDPAddr("udp", ":12121")
	if err != nil {
		log.Fatal(err.Error())
	}
	// targetU, err := net.ResolveUDPAddr("udp", "192.168.1.35:31861")

	udpListen, err := net.ListenUDP("udp", u)
	if err != nil {
		log.Fatal(err.Error())
	}

	b := make([]byte, 65535)
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

	for {
		n, _, err := udpListen.ReadFrom(b)
		atomic.AddInt64(&count, int64(n))
		if err != nil {
			// log.Fatal(err.Error())
		}

		// udpconn.Write(b[:n])

	}

}
