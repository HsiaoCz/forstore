package main

import (
	"fmt"
	"io"
	"log"

	"github.com/HsiaoCz/forstore/p2p"
)

type FileServerOpts struct {
	ListenAddr        string
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
}

type FileServer struct {
	FileServerOpts

	store *Store

	quitch chan struct{}
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
	}
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
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
