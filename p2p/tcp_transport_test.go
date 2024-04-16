package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":4001"
	tr := NewTCPTransport(listenAddr)

	assert.Equal(t, tr.listenAddress, listenAddr)

	// server
	// tr.Start()
	assert.Nil(t, tr.ListenAndAccept())
}
