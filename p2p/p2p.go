package p2p

import (
	"github.com/jxofficial/torrent-go/client"
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
	c, err := client.New(peer, t.PeerID, t.InfoHash)
	if err != nil {
		log.Printf("Could not perform handshake with %s. Disconnecting.\n", peer.IP)
	}
	defer c.Conn.Close()
	log.Printf("Completed handshake with %s\n", peer.IP)


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
		// start a goroutine for each peer (worker)
		go t.startDownloadWorker(peer, workQueue, results)
	}

	buf := make([]byte, t.Length)
	donePieces := 0
	// while there are still pieces left to process
	// continue reading from the results channel
	for donePieces < len(t.PieceHashes) {
		res := <-results
		// find byte at which piece starts
		begin, end := t.calculateBoundsForPiece(res.index)
		copy(buf[begin:end], res.buf)
		donePieces++

		percent := float64(donePieces) / float64(len(t.PieceHashes)) * 100
		log.Printf("(%0.2f%%) Downloaded. Piece %d", percent, res.index)
	}

	close(workQueue)
	return buf, nil
}