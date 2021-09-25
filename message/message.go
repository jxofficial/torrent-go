package message

import (
	"encoding/binary"
	"fmt"
	"io"
)

type messageID uint8

const (
	// MsgChoke chokes the receiver
	MsgChoke messageID = 0
	// MsgUnchoke unchokes the receiver
	MsgUnchoke messageID = 1
	// MsgInterested expresses interest in receiving data
	MsgInterested messageID = 2
	// MsgNotInterested expresses disinterest in receiving data
	MsgNotInterested messageID = 3
	// MsgHave alerts the receiver that the sender has downloaded a piece
	MsgHave messageID = 4
	// MsgBitfield encodes which pieces that the sender has downloaded
	MsgBitfield messageID = 5
	// MsgRequest requests a block of data from the receiver
	MsgRequest messageID = 6
	// MsgPiece delivers a block of data to fulfill a request
	MsgPiece messageID = 7
	// MsgCancel cancels a request
	MsgCancel messageID = 8
)

type Message struct {
	ID messageID
	Payload []byte
}

// FormatRequest creates a REQUEST message
func FormatRequest(index, begin, length int) *Message {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))
	return &Message{ID: MsgRequest, Payload: payload}
}

// FormatHave creates a HAVE message
func FormatHave(index int) *Message {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, uint32(index))
	return &Message{ID: MsgHave, Payload: payload}
}

// Serialize serializes a message struct into a buffer of the form
// <ID length+payload length><messageID><payload>
func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4) // message length is 0
	}
	len := uint32(len(m.Payload)+1) // length of message + length of id

	buf := make([]byte, 4+len)
	binary.BigEndian.PutUint32(buf[0:4], len)
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)
	return buf
}

// Read parses a message from a stream. Returns nil on keep-alive message
func Read(r io.Reader) (*Message, error) {
	lenBuf := make([]byte, 4)
	// r.read(buf) reads up to len(buf) bytes from r into buf
	_, err := io.ReadFull(r, lenBuf)
	if err != nil {
		return nil, err
	}
	len := binary.BigEndian.Uint32(lenBuf)

	// keep-alive message
	if len == 0 {
		return nil, nil
	}

	messageBuf := make([]byte, len)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil {
		return nil, err
	}

	m := &Message{
		ID: messageID(messageBuf[0]),
		Payload: messageBuf[1:],
	}
	return m, nil
}

func (m *Message) name() string {
	if m == nil {
		return "KeepAlive"
	}
	switch m.ID {
	case MsgChoke:
		return "Choke"
	case MsgUnchoke:
		return "Unchoke"
	case MsgInterested:
		return "Interested"
	case MsgNotInterested:
		return "NotInterested"
	case MsgHave:
		return "Have"
	case MsgBitfield:
		return "Bitfield"
	case MsgRequest:
		return "Request"
	case MsgPiece:
		return "Piece"
	case MsgCancel:
		return "Cancel"
	default:
		return fmt.Sprintf("Unknown#%d", m.ID)
	}
}

func (m *Message) String() string {
	if m == nil {
		return m.name()
	}
	return fmt.Sprintf("%s [%d]", m.name(), len(m.Payload))
}