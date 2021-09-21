package p2p

import (
	"github.com/jxofficial/torrent-go/peers"
	"log"
)

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

type pieceResult struct {
	index int
	buf []byte
}

// index is the index of hashedPiece
// ie the piece index number
// since each piece has fixed length, we can calculate where the bytes
// for a specific piece starts
// end is exclusive
func (t *Torrent) calculateBoundsForPiece(index int) (begin int, end int) {
	begin = index * t.PieceLength
	end = begin + t.PieceLength
	if end > t.Length {
		end = t.Length
	}
	return begin, end
}

func (t *Torrent) calculatePieceSize(index int) int {
	begin, end := t.calculateBoundsForPiece(index)
	return end - begin
}

func (t *Torrent) startDownloadWorker(
	peer peers.Peer, workQueue chan *pieceWork, results chan *pieceResult) {
	// c, err := client.New(peer, t.PeerID, t.InfoHash)
}


func (t *Torrent) Download() ([]byte, error) {
	log.Println("Starting download for", t.Name)
	// each piece has one corresponding piece hash
	workQueue := make(chan *pieceWork, len(t.PieceHashes))
	results := make(chan *pieceResult)

	for index, hash := range t.PieceHashes {
		length := t.calculatePieceSize(index)
		// work queue contains all the pieces that will need to be picked up by the peers
		// index is the index of the piece
		workQueue <- &pieceWork{index, hash, length}
	}

	for _, peer := range t.Peers {
		// start a goroutine for each peer
		go t.startDownloadWorker(peer, workQueue, results)
	}

	return make([]byte, 0), nil
}