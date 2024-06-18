package main

import (
	"fmt"
	"os"
	"log"
	"github.com/Triplem5ds/BittorrentClient/torrentfile"
	
)
func main() {
	inputFilePath := os.Args[1]
	// outputFilePath := os.Args[2]

	torrentFile, err := torrentfile.Open(inputFilePath)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v", torrentFile)


}
