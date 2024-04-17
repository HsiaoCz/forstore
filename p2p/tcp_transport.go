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

type TCPTransportOps struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOps
	listener net.Listener
	mu       sync.RWMutex
	peers    map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOps) *TCPTransport {
	return &TCPTransport{
		TCPTransportOps: opts,
	}
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
	peer := NewTCPPeer(conn, true)
	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		slog.Error("TCP handshake error", "err", err)
		return
	}

	// read loop
	// msg := &Message{}
	buf := make([]byte, 2000)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			slog.Error("TCP read error", "err", err)
		}
		// if err := t.Decoder.Decode(conn, buf); err != nil {
		// 	slog.Error("conn read error", "err", err)
		// 	continue
		// }
		fmt.Printf("message:%+v\n", buf[:n])
	}
}
