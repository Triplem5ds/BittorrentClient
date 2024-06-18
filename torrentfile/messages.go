package torrentfile

import (
	"fmt"
	"bytes"
	"crypto/sha1"

	"github.com/jackpal/bencode-go"
)

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

func (i* bencodeInfo) hash() ([20]byte, error) {
	///Hash the info to check credibelity 

	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (i* bencodeInfo) splitPieceHashes() ([][20]byte, error) {
	///take the piece hashes and split it into an array for pieces
	///where ever arr[i] = hashPiece[i]
	buf := []byte(i.Pieces)
	if len(buf) % 20 !=0 {
		return nil, fmt.Errorf("Received malformed pieces of length %d", len(buf))
	}

	hashesNum := len(buf) / 20
	hashes := make([][20]byte, hashesNum)

	for i := 0; i < hashesNum; i++ {
		copy(hashes[i][:], buf[i*20:(i + 1) * 20])
	}

	return hashes, nil
}

func (bto *bencodeTorrent) toTorrentfile() (TorrentFile, error) {
	///	transforms the bencode torrent to a torrent file structure

	infoHash, err := bto.Info.hash()

	if err != nil {
		return TorrentFile{}, err
	}

	piecesHash, err := bto.Info.splitPieceHashes()
	if err != nil {
		return TorrentFile{}, err
	}
	
	return TorrentFile {
		Announce: bto.Announce,
		InfoHash: infoHash,
		PieceHashes: piecesHash,
		PieceLength: bto.Info.PieceLength,
		Length: bto.Info.Length,
		Name: bto.Info.Name,
	}, nil

}