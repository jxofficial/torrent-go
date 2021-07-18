package handshake

import (
	"fmt"
	"io"
)

// A Handshake is a special message that a peer uses to identify itself
type Handshake struct {
	Pstr string // protocol identifier, which is always "BitTorrent protocol"
	InfoHash [20]byte
	PeerID [20]byte
}

// New creates a new handshake
func New(infoHash, peerID [20]byte) *Handshake {
	return &Handshake{
		Pstr: "BitTorrent protocol",
		InfoHash: infoHash,
		PeerID: peerID,
	}
}

// Serialize converts the handshake into a buffer
func (h *Handshake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr)+49)
	buf[0] = byte(len(h.Pstr))
	curr := 1
	curr += copy(buf[curr:], h.Pstr)
	curr += copy(buf[curr:], make([]byte, 8))
	curr += copy(buf[curr:], h.InfoHash[:])
	curr += copy(buf[curr:], h.PeerID[:])
	return buf
}

func Read(r io.Reader) (*Handshake, error) {
	// remove the first byte which is the pstr length
	bufLen := make([]byte, 1)
	_, err := io.ReadFull(r, bufLen)
	if err != nil {
		return nil, err
	}
	pstrLen := int(bufLen[0])

	if pstrLen == 0 {
		err := fmt.Errorf("protocol string (pstr) cannot be of length 0")
		return nil, err
	}

	// 20 bytes for infoHash and peerID respectively
	// 8 bytes for extensions flag
	handshakeBuf := make([]byte, 48+pstrLen)
	_, err = io.ReadFull(r, handshakeBuf)
	if err != nil {
		return nil, err
	}

	var infoHash, peerID [20]byte
	copy(infoHash[:], handshakeBuf[pstrLen+8:pstrLen+28])
	copy(peerID[:], handshakeBuf[pstrLen+28:])

	h := &Handshake{
		Pstr:     string(handshakeBuf[0:pstrLen]),
		InfoHash: infoHash,
		PeerID:   peerID,
	}
	return h, nil
}
