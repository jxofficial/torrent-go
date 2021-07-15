package peers

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

type Peer struct {
	IP net.IP
	Port uint16
}

func Unmarshal(peersBinary []byte) ([]Peer, error) {
	const peerBytes = 6
	if len(peersBinary) % peerBytes != 0 {
		err := fmt.Errorf("Received malformed peers")
		return nil, err
	}

	numPeers := len(peersBinary) / peerBytes
	peers := make([]Peer, numPeers)
	for i := 0; i < numPeers; i++ {
		offset := i * peerBytes
		peers[i].IP = peersBinary[offset : offset+4]
		peers[i].Port = binary.BigEndian.Uint16(peersBinary[offset+4 : offset+6])
	}

	return peers, nil
}

func (p Peer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}