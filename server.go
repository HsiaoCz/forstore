package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"log/slog"
	"sync"

	"github.com/HsiaoCz/forstore/p2p"
)

type FileServerOpts struct {
	ListenAddr        string
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootStrapNodes    []string
}

type FileServer struct {
	FileServerOpts

	store *Store

	peerLock sync.Mutex

	quitch chan struct{}
	peers  map[string]p2p.Peer
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

type Payload struct {
	Key  string
	Data []byte
}

func (s *FileServer) broadcast(p Payload) error {
	// buf := new(bytes.Buffer)
	// for _, peer := range s.peers {
	// 	if err := gob.NewEncoder(buf).Encode(p); err != nil {
	// 		return err
	// 	}
	// 	peer.Send(buf.Bytes())
	// }
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(p)
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	// 1. store the file to disk
	// 2. broadcast this file to all konwn peers in the network

	if err := s.store.Write(key, r); err != nil {
		return err
	}

	// the reader is empty
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, r)
	if err != nil {
		return err
	}

	fmt.Println(buf.Bytes())

	return nil
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p

	log.Printf("connected with remote %s", p.RemoteAddr())

	return nil
}

func (s *FileServer) bootStrapNetwork() error {
	for _, addr := range s.BootStrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			if err := s.Transport.Dial(addr); err != nil {
				slog.Error("dial error", "err", err)
			}
		}(addr)
	}
	return nil
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	s.bootStrapNetwork()
	s.loop()
	return nil
}

func (s *FileServer) loop() {

	defer func() {
		log.Println("file server stopped user quit action")
		if err := s.Transport.Close(); err != nil {
			log.Fatalf("tcp transport close failed %v\n", err)
		}
	}()

	for {
		select {
		case msg := <-s.Transport.Consume():
			fmt.Println(msg)
		case <-s.quitch:
			return
		}
	}
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) Store(key string, r io.Reader) error {
	return s.store.Write(key, r)
}
