package p2p

import "errors"

// ErrInvalidHandshake is returned if the handshake between
// the local and reomte node could not be established.
var ErrInvalidHandshake = errors.New("invalid handshake")

// HandshakeFunc ...?
type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(Peer) error { return nil }

type OnPeerFunc func(Peer) error

func NOPOnPeerFunc(Peer) error { return errors.New("failed the onpeer func") }
