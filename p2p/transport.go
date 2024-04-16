package p2p

// Peer is an interface that  repersents the remote node
type Peer interface {
}

// Transport is anything that handles the communication
// between the nodes in the network
// form (TCP,UDP,websockets or any others...)
type Transport interface {
	ListenAndAccept() error
}
