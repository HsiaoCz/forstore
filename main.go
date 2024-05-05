package main

import (
	"bytes"
	"log"
	"time"

	"github.com/HsiaoCz/forstore/p2p"
)

// func main() {
// 	tcpOpts := p2p.TCPTransportOps{
// 		ListenAddr:    ":3001",
// 		HandshakeFunc: p2p.NOPHandshakeFunc,
// 		Decoder:       p2p.DefaultDecoder{},
// 		OnPeer:        p2p.NOPOnPeerFunc,
// 	}
// 	tr := p2p.NewTCPTransport(tcpOpts)

// 	go func() {
// 		for {
// 			msg := <-tr.Consume()
// 			fmt.Printf("%+v\n", msg)
// 		}
// 	}()

// 	if err := tr.ListenAndAccept(); err != nil {
// 		log.Fatal(err)
// 	}

// 	select {}
// }

func main() {
	s1 := makeServer(":3003", "")
	s2 := makeServer(":4000", ":3003")
	go func() {
		// time.Sleep(time.Second * 1)
		log.Fatal(s1.Start())
	}()

	time.Sleep(time.Second)

	go s2.Start()

	// time.Sleep(time.Second * 1)

	data := bytes.NewReader([]byte("my big data file here"))

	s2.StoreData("key", data)

	select {}
}

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOps{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// TODO onPeer func
		// OnPeer:        p2p.NOPOnPeerFunc,
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)
	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootStrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer

	return s
}
