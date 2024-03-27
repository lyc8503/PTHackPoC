package main

import (
	"compress/gzip"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/lyc8503/ptcheat/util"
)

var port int

func requestTracker(trackerUrl string, infoHashHex string, initSize int64) {
	infoHash, err := hex.DecodeString(infoHashHex)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	infoHashEncoded := url.QueryEscape(string(infoHash))

	peerId := util.RandomPeerId()
	key := util.RandomKey()

	reqURL := fmt.Sprintf("%s&info_hash=%s&peer_id=%s&port=%d&uploaded=0&downloaded=0&left=%d&corrupt=0&key=%s"+
		"&event=started&numwant=200&compact=1&no_peer_id=1&supportcrypto=1&redundant=0",
		trackerUrl, infoHashEncoded, peerId, port, initSize, key)

	fmt.Println("Requesting:", reqURL)

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", "qBittorrent/4.6.3")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Connection", "close")

	// TODO: fake TLS handshake fingerprint
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	body, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	filename := fmt.Sprintf("%s.peers", infoHashHex)
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	_, err = file.Write(body)
	if err != nil {
		panic(err)
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}

	// fmt.Println(string(body))

	// var v interface{}
	// err = bencode.Unmarshal(body, &v)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("=====\n%+v\n", v)
}

func localFakeTrackerHandler(w http.ResponseWriter, r *http.Request) {
	infoHash := r.URL.Query().Get("info_hash")
	infoHash = hex.EncodeToString([]byte(infoHash))
	fmt.Println("received request: ", infoHash)

	filename := fmt.Sprintf("%s.peers", infoHash)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("file open error: ", err)
		return
	}
	cachedPeers, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("file read error: ", err)
		return
	}

	w.Write(cachedPeers)
}

func main() {
	// spam peer
	if len(os.Args) >= 2 {
		util.ConnectPeer(os.Args[1], os.Args[2])
		return
	}

	// Below is the main logic to process *.torrent files
	// TODO: maybe generate a fixed port number on first run
	// ideally use a port number that is identical to the one in your BT client
	port = util.RandomPort()

	// List *.torrent files
	files, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(files); i++ {
		if strings.HasSuffix(files[i].Name(), ".torrent") {
			if strings.HasPrefix(files[i].Name(), "FREE_") {
				fmt.Println("skipping already processed torrent: ", files[i].Name())
				continue
			}

			fmt.Println("processing: ", files[i].Name())
			realAnnounce, hash, leftSize, err := util.ParseAndRegenerateTorrent(files[i].Name(), "http://127.0.0.1:1088/announce")
			if err != nil {
				fmt.Println("Error: ", err)
				continue
			}
			fmt.Printf("info_hash: %s, size: %d\n", hash, leftSize)

			requestTracker(
				realAnnounce,
				hash,
				leftSize,
			)
		}
	}

	http.HandleFunc("/announce", localFakeTrackerHandler)

	fmt.Println("Starting server at port 1088")
	if err := http.ListenAndServe(":1088", nil); err != nil {
		panic(err)
	}
}
