package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/quic-go/quic-go"
	. "github.com/someview/quic-benchmark"
)

var sendCount = int64(0)
var recvCount = int64(0)
var maxClientNum = 10

var multiMode = 0  // 大量客户端，均发送消息
var singleMode = 1 // 大量客户端，只有一个客户端在发送消息
var slientMode = 2 // 大量客户端，不发送消息

func RunClient(runMode int) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	conn, err := quic.DialAddr(ctx, Addr, tlsConf, QuicConf)
	if err != nil {
		fmt.Println("dial err:", err)
		return
	}

	if runMode != multiMode {
		return
	}

	go func() {
		maxData := []byte("hello 123456789")
		var err error
		time.Sleep(time.Second * 20)
		for {
			if err = conn.SendDatagram(maxData); err != nil {
				fmt.Println("cli send dgram error:", err)
				return
			}
			atomic.AddInt64(&sendCount, 1)
		}
	}()

	go func() {
		// for range time.NewTicker(time.Microsecond).C {
		//maxData := make([]byte, 15)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		for {
			if recvData, err := conn.ReceiveDatagram(ctx); err != nil {
				log.Fatalln("recvData error:", recvData)
			}
			atomic.AddInt64(&recvCount, 1)
		}
	}()
}

func ReportView() {
	for range time.NewTicker(time.Second * 20).C {
		send := atomic.LoadInt64(&sendCount)
		recv := atomic.LoadInt64(&recvCount)
		fmt.Println("时间:", time.Now(), "send:", send, "recv:", recv, "rate:", (send+recv)/20)
		atomic.StoreInt64(&sendCount, 0)
		atomic.StoreInt64(&recvCount, 0)
	}
}

func main() {
	fmt.Println("开始时间:", time.Now())
	var mode = flag.Int("mode", 0, "运行模式")
	flag.Parse()
	durtion := time.Microsecond * 500
	timer := time.NewTimer(durtion)
	switch *mode {
	case slientMode:
		for i := 0; i < int(maxClientNum); i++ {
			<-timer.C
			go RunClient(slientMode)
			timer.Reset(durtion)
		}
	case singleMode:
		for i := 0; i < int(maxClientNum)-1; i++ {
			<-timer.C
			go RunClient(slientMode)
			timer.Reset(durtion)
		}
		go RunClient(multiMode)
	case multiMode:
		for i := 0; i < int(maxClientNum); i++ {
			<-timer.C
			go RunClient(multiMode)
			timer.Reset(durtion)
		}
	}
	go ReportView()
	time.Sleep(time.Minute * 30)
}
