package p2p

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection
type TCPPeer struct {
	// conn is the underlying connection of the peer
	net.Conn
	// if we dial and retrieve a conn => outbound ==true
	// if we accept and retrieve a conn => outbound ==false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outBound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outBound,
	}
}

// Close implements the tcp networks
// func (p *TCPPeer) Close() error {
// 	return p.conn.Close()
// }

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
}

// func (p *TCPPeer) Conn() net.Conn {
// 	return p.conn
// }

// RemoteAddr implements the Peer interface
// and will return the remote address of
// its underlying connection
// func (p *TCPPeer) RemoteAddr() net.Addr {
// 	return p.conn.RemoteAddr()
// }

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

// Close implements the Transport interface

func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Consume implements the Transport interface,which will return read-only channel
// for reading the incoming messages received from another peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

// Dial implements the Transport interface .
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)

	return nil
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()

	log.Printf("TCP transport listening on port: %s\n", t.ListenAddr)

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			slog.Error("TCP accept error", "err", err)
		}

		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		slog.Error("dropping peer connection", "err", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound)
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
		err := t.Decoder.Decode(conn, &rpc)
		if _, ok := err.(*net.OpError); ok {
			return
		}
		if err != nil {
			slog.Error("tcp connection read error", "err", err)
			continue
		}
		rpc.From = conn.RemoteAddr()
		fmt.Printf("message:%+v\n", rpc)
	}
}
