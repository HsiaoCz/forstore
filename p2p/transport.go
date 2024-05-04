package p2p

import "net"

// Peer is an interface that  repersents the remote node
type Peer interface {
	// Conn() net.Conn
	Send([]byte) error
	// RemoteAddr() net.Addr
	// Close() error

	// the underlying connection of the peer. whick in this case is a TCP connection
	net.Conn
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
