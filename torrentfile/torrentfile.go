package torrentfile

import (
	"os"

	"github.com/jackpal/bencode-go"
)

func Open(path string) (TorrentFile, error) {

	file, err := os.Open(path)

	if err != nil {
		return TorrentFile{}, err
	}
	
	defer file.Close()	///don't keep files open because of memory & resources

	bto := bencodeTorrent{}
	err = bencode.Unmarshal(file, &bto)

	if err != nil {
		return TorrentFile{}, err
	}

	return bto.toTorrentfile()
}