package torrentfile

import (
	"os"

	"github.com/zeebo/bencode"
)

func Open(path string) (TorrentFile, error) {

	file, err := os.ReadFile(path)

	if err != nil {
		return TorrentFile{}, err
	}

	bto := bencodeTorrent{}
	err = bencode.DecodeBytes(file, &bto)

	if err != nil {
		return TorrentFile{}, err
	}

	if err != nil {
		return TorrentFile{}, err
	}

	return bto.toTorrentfile()
}
