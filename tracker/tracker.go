package tracker

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/zeebo/bencode"
)

const Port uint16 = 6881

type TrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func getAddressesFromInitalResponse(resp initialHTTPTrackerResponse) ([]net.TCPAddr, error) {

	if len(resp.Peers)%6 != 0 {
		return nil, fmt.Errorf("malformed http tracker response")
	}

	res := make([]net.TCPAddr, len(resp.Peers)/6)

	for i := 0; i < len(resp.Peers)/6; i++ {
		rawPort := []byte(resp.Peers[i*6+4 : i*6+6])
		port := binary.BigEndian.Uint16(rawPort)

		res = append(res, net.TCPAddr{
			IP:   []byte(resp.Peers[i*6 : i*6+4]),
			Port: int(port),
		})
	}

	return res, nil
}

func getAddressesFromExpandedResponse(resp expandedHTTPTrackerResponse) ([]net.TCPAddr, error) {

	res := make([]net.TCPAddr, len(resp.Peers))

	for _, peer := range resp.Peers {
		res = append(res, net.TCPAddr{
			IP:   net.ParseIP(peer.IP),
			Port: peer.Port,
		})
	}

	return res, nil
}

func getPeersFromHTTPTracker(u *url.URL, infoHash, peerId [20]byte) ([]net.TCPAddr, error) {
	///returns a list of peers ips
	params := url.Values{
		"info_hash":  []string{string(infoHash[:])},
		"peer_id":    []string{string(peerId[:])},
		"port":       []string{strconv.Itoa(int(Port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{"0"},
	}

	u.RawQuery = params.Encode()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 response %d %s", resp.StatusCode, string(raw))
	}

	var initialResponse initialHTTPTrackerResponse
	var tmp interface{}

	err = bencode.DecodeBytes(raw, &tmp)
	err = bencode.DecodeBytes(raw, &initialResponse)
	if err == nil { /// The response is in the initial format
		return getAddressesFromInitalResponse(initialResponse)
	}

	var expandedResponse expandedHTTPTrackerResponse
	err = bencode.DecodeBytes(raw, &expandedResponse)

	if err != nil {
		return nil, err
	}

	return getAddressesFromExpandedResponse(expandedResponse)
}

func GetPeers(link string, infoHash, peerId [20]byte) ([]net.TCPAddr, error) {
	///build the tracker URL to get the peers
	base, err := url.Parse(link)

	if err != nil {
		return nil, err
	}

	switch base.Scheme {
	case "http", "https":
		return getPeersFromHTTPTracker(base, infoHash, peerId)
	default:
		return nil, fmt.Errorf("Unrecognized url scheme %s", base.Scheme)
	}
}
