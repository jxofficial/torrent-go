package torrentfile

import (
	"github.com/jackpal/bencode-go"
	"github.com/jxofficial/torrent-go/peers"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type bencodeTrackerResp struct {
	Interval int `bencode:"interval"`
	// Peers is binary blob of all the IP addresses of each peer, grouped by 6 bytes each
	// 4 bytes for IP, 2 bytes for port represented as a big-endian uint16
	Peers string `bencode:"peers"`
}


func (t *TorrentFile) buildTrackerURL(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values {
		"info_hash": []string{string(t.InfoHash[:])}, // identifies the file we wish to download
		"peer_id": []string{string(peerID[:])}, // random id for self identification
		"port": []string{strconv.Itoa(int(port))},
		"uploaded": []string{"0"},
		"downloaded": []string{"0"},
		"compact": []string{"1"},
		"left": []string{strconv.Itoa(t.Length)},
	}
	base.RawQuery = params.Encode()
	// String is a method on the *URL receiver
	return base.String(), nil
}

func (t *TorrentFile) requestPeers(peerID [20]byte, port uint16) ([]peers.Peer, error) {
	trackerURL, err := t.buildTrackerURL(peerID, port)
	if err != nil {
		return nil, err
	}

	c := &http.Client{Timeout: 15 * time.Second}
	resp, err := c.Get(trackerURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	trackerResp := bencodeTrackerResp{}
	err = bencode.Unmarshal(resp.Body, &bencodeTrackerResp{})
	if err != nil {
		return nil, err
	}

	return peers.Unmarshal([]byte(trackerResp.Peers))
}