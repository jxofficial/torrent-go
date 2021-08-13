package p2p

import "github.com/jxofficial/torrent-go/peers"

// Torrent holds data required to download a torrent from a list of peers
type Torrent struct {
	Peers []peers.Peer
	PeerID [20]byte
	InfoHash [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length int
	Name string
}

type pieceWork struct {
	index int
	hash [20]byte
	length int
}

