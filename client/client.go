package client

import (
	"bytes"
	"fmt"
	"github.com/jxofficial/torrent-go/handshake"
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

// completeHandshake sends a handshake to the peer
// and checks if the peer's handshake response contains the same infoHash
func completeHandshake(conn net.Conn, infoHash, peerID [20]byte) (*handshake.Handshake, error) {
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{}) // disable the deadline

	req := handshake.New(infoHash, peerID)
	_, err := conn.Write(req.Serialize())
	if err != nil {
		return nil, err
	}

	resp, err := handshake.Read(conn)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(resp.InfoHash[:], infoHash[:]) {
		return nil, fmt.Errorf("expected infohash to be %x but got %x",
			infoHash, resp.InfoHash)
	}

	return resp, nil
}

// New connects with a peer and completes a handshake
// returns an error if any of these steps fail
func New(peer peers.Peer, peerID, infoHash [20]byte) (*Client, error) 	{
	// TODO: implement logic
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
	if err != nil {
		return nil, err
	}

	_, err = completeHandshake(conn, infoHash, peerID)



	return &Client{}, nil
}
