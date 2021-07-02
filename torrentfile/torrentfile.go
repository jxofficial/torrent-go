package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/jackpal/bencode-go"
	"os"
)

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"` // binary blob of the hashes of each piece
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

// TorrentFile is an application layer struct
type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte   // hash of bencode info - including pieces, length, name, piece length
	PieceHashes [][20]byte // slice of 20byte arrays
	PieceLength int
	Length      int
	Name        string
}

// Open parses a torrent file
// it de-serializes the torrent file into a bencodeTorrent serialization struct
// the bencodeTorrent serialization struct is then converted to a TorrentFile application struct
func Open(path string) (TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, nil
	}
	defer file.Close()

	bto := bencodeTorrent{}
	// reads bencode data and transforms in into a bencodeTorrent struct
	err = bencode.Unmarshal(file, &bto)
	if err != nil {
		return TorrentFile{}, err
	}
	return bto.toTorrentFile()
}

// splitPiecesHashes splits the hashes string into individual hash of each piece
// each hash is 20 bytes
func (info *bencodeInfo) splitPiecesHashes() ([][20]byte, error) {
	const hashLen = 20
	buffer := []byte(info.Pieces) // convert string into byte slice
	if len(buffer) % hashLen != 0 {
		err := fmt.Errorf("received malformed pieces of length %d", len(buffer))
		return nil, err
	}
	numHashes := len(buffer) / hashLen
	hashes := make([][20]byte, 0, numHashes)

	for i := 0; i < numHashes; i++ {
		var hashArr [20]byte
		copy(hashArr[:], buffer[i*hashLen:(i+1)*hashLen])
		hashes = append(hashes, hashArr)
	}
	return hashes, nil
}

// hash calculates the checksum of the bencodeInfo struct
func (info *bencodeInfo) hash() ([20]byte, error) {
	var buffer bytes.Buffer
	err := bencode.Marshal(&buffer, *info)
	if err != nil {
		return [20]byte{}, err
	}
	hash := sha1.Sum(buffer.Bytes())
	return hash, nil
}

// converts a bencodeTorrent struct into a TorrentFile application struct
func (bto bencodeTorrent) toTorrentFile() (TorrentFile, error) {
	infoHash, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}
	pieceHashes, err := bto.Info.splitPiecesHashes()

	t := TorrentFile{
		Announce:    bto.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}
	return t, nil
}

