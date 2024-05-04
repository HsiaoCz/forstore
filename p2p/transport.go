package p2p

import "net"

// Peer is an interface that  repersents the remote node
type Peer interface {
	RemoteAddr() net.Addr
	Close() error
}

// Transport is anything that handles the communication
// between the nodes in the network
// form (TCP,UDP,websockets or any others...)
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
	Dial(string) error
}
