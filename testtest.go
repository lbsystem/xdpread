package main

import (
	"encoding/binary"
	"fmt"

	"time"
)

func main() {
	now:=time.Now()
	var b int64
	for i := 0; i < 1000000000; i++ {
		b=int64(binary.BigEndian.Uint16([]byte{122,33}))
	}
	fmt.Println(time.Since(now))
	fmt.Println(b)
}