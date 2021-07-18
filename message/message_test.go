package message

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerialize(t *testing.T) {
	tests := map[string]struct {
		input  *Message
		output []byte
	}{
		"serialize message": {
			input:  &Message{ID: MsgHave, Payload: []byte{1, 2, 3, 4}},
			output: []byte{0, 0, 0, 5, 4, 1, 2, 3, 4},
		},
		"serialize keep-alive": {
			input:  nil,
			output: []byte{0, 0, 0, 0},
		},
	}

	for _, test := range tests {
		buf := test.input.Serialize()
		assert.Equal(t, test.output, buf)
	}
}

func TestRead(t *testing.T) {
	tests := map[string]struct{
		input []byte
		output *Message
		fails bool
	}{
		"parse normal message into struct": {
			input: []byte{0, 0, 0, 5, 4, 1, 2, 3, 4},
			output: &Message{
				ID: 4,
				Payload: []byte{1, 2, 3, 4},
			},
			fails: false,
		},
		"parse keep alive into nil": {
			input: []byte{0, 0, 0, 0},
			output: nil,
			fails: false,
		},
		"length too short": {
			input: []byte{0, 0, 1},
			output: nil,
			fails: true,
		},
		"buffer not equals to length": {
			input:  []byte{0, 0, 0, 5, 4, 1, 2},
			output: nil,
			fails:  true,
		},
	}

	for _, test := range tests {
		// this reader implements the io.Reader interface
		// ie it has a concrete implementation of a
		// Read(buf []byte) (n int, err error) method
		reader := bytes.NewReader(test.input)
		// Read internally calls io.ReadFull(reader, destinationBuffer)
		// reader is alr initialized with a buffer to read from, in this case test.input
		m, err := Read(reader)
		if test.fails {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
		assert.Equal(t, m, test.output)
	}
}