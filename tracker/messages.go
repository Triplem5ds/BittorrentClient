package tracker

type initialHTTPTrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

type expandedHTTPTrackerResponse struct {
	Peers []struct {
		ID   string `bencode:"peer_id"`
		IP   string `bencode:"ip"`
		Port int    `bencode:"port"`
	} `bencode:"peers"`
	// a lot of fields that are not used
	Interval   int    `bencode:"interval"`  // likely needs to be escaped
	InfoHash   string `bencode:"info_hash"` // likely needs to be escaped
	Uploaded   int    `bencode:"uploaded"`
	Downloaded int    `bencode:"downloaded"`
	Left       int    `bencode:"left"`
	Event      string `bencode:"event"`
}
