package p2p

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection
type TCPPeer struct {
	// conn is the underlying connection of the peer
	conn net.Conn
	// if we dial and retrieve a conn => outbound ==true
	// if we accept and retrieve a conn => outbound ==false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outBound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outBound,
	}
}

// Close implements the tcp networks
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOps struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        OnPeerFunc
}

type TCPTransport struct {
	TCPTransportOps
	listener net.Listener
	rpcch    chan RPC
	mu       sync.RWMutex
	peers    map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOps) *TCPTransport {
	return &TCPTransport{
		TCPTransportOps: opts,
		rpcch:           make(chan RPC),
	}
}

// Consume implements the Transport interface,which will return read-only channel
// for reading the incoming messages received from another peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			slog.Error("TCP accept error", "err", err)
		}

		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error

	defer func() {
		slog.Error("dropping peer connection", "err", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)
	if err = t.HandshakeFunc(peer); err != nil {
		slog.Error("TCP handshake error", "err", err)
		return
	}
	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			slog.Error("on peer error", "err", err)
			return
		}
	}
	// read loop
	rpc := RPC{}
	// buf := make([]byte, 2000)
	for {
		// n, err := conn.Read(buf)
		// if err != nil {
		// 	slog.Error("TCP read error", "err", err)
		// }
		if err = t.Decoder.Decode(conn, &rpc); err != nil {
			slog.Error("conn read error", "err", err)
			continue
		}
		rpc.From = conn.RemoteAddr()
		fmt.Printf("message:%+v\n", rpc)
	}
}
