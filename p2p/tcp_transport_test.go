package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	tcpOpts := TCPTransportOps{
		ListenAddr:    ":3001",
		HandshakeFunc: NOPHandshakeFunc,
		Decoder:       DefaultDecoder{},
	}
	tr := NewTCPTransport(tcpOpts)

	assert.Equal(t, tr.ListenAddr, tcpOpts.ListenAddr)

	// server
	// tr.Start()
	assert.Nil(t, tr.ListenAndAccept())
}
