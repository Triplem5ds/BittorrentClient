package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/Triplem5ds/BittorrentClient/torrentfile"
	"github.com/Triplem5ds/BittorrentClient/tracker"
)

func dedupeAddrs(addrs []net.TCPAddr) []net.TCPAddr {
	deduped := []net.TCPAddr{}
	set := map[string]bool{}
	for _, a := range addrs {
		if !set[a.String()] {
			deduped = append(deduped, a)
			set[a.String()] = true
		}
	}
	return deduped
}

func main() {

	// inputFilePath := os.Args[1]
	// outputFilePath := os.Args[2]
	inputFilePath := "testData/ubuntu-20.04.2-live-server-amd64.iso.torrent"

	torrentFile, err := torrentfile.Open(inputFilePath)

	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	var mut sync.Mutex
	wg.Add(len(torrentFile.URLs))
	var peerAddrs []net.TCPAddr

	var peerId [20]byte
	_, err = rand.Read(peerId[:])

	if err != nil {
		log.Fatal(err)
	}

	for _, link := range torrentFile.URLs {
		go func() {
			defer wg.Done()
			result, err := tracker.GetPeers(link, torrentFile.InfoHash, [20]byte(peerId))
			if err != nil {
				fmt.Printf("error getting peers from tracker %s: %s\n", link, err.Error())
				return
			}
			mut.Lock()
			fmt.Printf("Got Peers from tracker %s: %d\n", link, len(result))
			peerAddrs = append(peerAddrs, result...)
			mut.Unlock()
		}()
	}

	wg.Wait()
	peerAddrs = dedupeAddrs(peerAddrs)

}
