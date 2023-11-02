package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"runtime"
	"sync/atomic"

	"github.com/quic-go/quic-go"
	. "github.com/someview/quic-benchmark"
)

var activeCount = int64(0)

// Start a server that echos all data on the first stream opened by the client
func echoServer() {
	listener, err := quic.ListenAddr(Addr, generateTLSConfig(), QuicConf)
	if err != nil {
		log.Fatalln("lisnten err:", err)
	}

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			fmt.Println("recv conn err:", err)
			continue
		}
		atomic.AddInt64(&activeCount, 1)
		fmt.Println("server recv quic conn:", atomic.LoadInt64(&activeCount), "routines:", runtime.NumGoroutine()) // 最终的活跃连接数
		go func() {
			for {
				stream, err := conn.AcceptStream(context.Background())
				if err != nil {
					fmt.Println("conn accept stream err:", err)
					return
				}
				if _, err := io.Copy(stream, stream); err != nil {
					fmt.Println("stream send err:", err)
					_ = stream.Close()
					return
				}
			}
		}()
	}
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}

// We start a server echoing data on the first stream the client opens,
// then connect with a client, send the message, and wait for its receipt.
func main() {
	echoServer()
	// time.Sleep(time.Minute * 30)
}
