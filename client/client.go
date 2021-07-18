package client

import (
	"github.com/jxofficial/torrent-go/peers"
	"net"
	"time"
)

// A Client is a TCP connection with a Peer
type Client struct {
	Conn     net.Conn
	Choked   bool // is the client choked by the peer
	// TODO: add field
	// Bitfield bitfield.Bitfield
	peer     peers.Peer
	infoHash [20]byte
	peerID   [20]byte
}

// New connects with a peer, completes and receives a handshake
// returns an err if any of these steps fail
func New(peer peers.Peer, peerID, infoHash [20]byte) (*Client, error) 	{
	// TODO: implement logic
	_, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
	if err != nil {
		return nil, err
	}
	return &Client{}, nil
}
