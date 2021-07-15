package torrentfile

import (
	"github.com/jxofficial/torrent-go/peers"
	"net/url"
	"strconv"
)

func (t *TorrentFile) buildTrackerURL(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values {
		"info_hash": []string{string(t.InfoHash[:])},
		"peer_id": []string{string(peerID[:])},
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
	return nil, nil
}