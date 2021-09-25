package client

import (
	"bytes"
	"fmt"
	"github.com/jxofficial/torrent-go/bitfield"
	"github.com/jxofficial/torrent-go/handshake"
	"github.com/jxofficial/torrent-go/message"
	"github.com/jxofficial/torrent-go/peers"
	"net"
	"time"
)

// A Client is a TCP connection with a Peer
type Client struct {
	Conn     net.Conn
	Choked   bool // is the client choked by the peer
	Bitfield bitfield.Bitfield
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

func recvBifield(conn net.Conn) (bitfield.Bitfield, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{}) // disable the deadline

	msg, err := message.Read(conn)
	if err != nil {
		return nil, err
	}

	if msg.ID != message.MsgBitfield {
		err := fmt.Errorf("expected bitfield by got ID %d", msg.ID)
		return nil, err
	}

	return msg.Payload, nil
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
	if err != nil {
		conn.Close()
		return nil, err
	}

	bf, err := recvBifield(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{
		Conn: conn,
		Choked: true,
		Bitfield: bf,
		peer: peer,
		infoHash: infoHash,
		peerID: peerID,
	}, nil
}


// Read reads and consumes a message from the connection
func (c *Client) Read() (*message.Message, error) {
	msg, err := message.Read(c.Conn)
	return msg, err
}

// SendRequest sends a Request message to the peer
func (c *Client) SendRequest(index, begin, length int) error {
	req := message.FormatRequest(index, begin, length)
	_, err := c.Conn.Write(req.Serialize())
	return err
}

// SendInterested sends an Interested message to the peer
func (c *Client) SendInterested() error {
	msg := message.Message{ID: message.MsgInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendNotInterested sends a NotInterested message to the peer
func (c *Client) SendNotInterested() error {
	msg := message.Message{ID: message.MsgNotInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendUnchoke sends an Unchoke message to the peer
func (c *Client) SendUnchoke() error {
	msg := message.Message{ID: message.MsgUnchoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendHave sends a Have message to the peer
func (c *Client) SendHave(index int) error {
	msg := message.FormatHave(index)
	_, err := c.Conn.Write(msg.Serialize())
	return err
}