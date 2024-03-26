package util

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"math/rand"
	"os"

	"github.com/anacrolix/torrent/bencode"
)

// Generate a random peer ID (qBit 4.6.3)
func RandomPeerId() string {
	chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_.!~*()"
	ret := "-qB4630-"

	for i := 0; i < 12; i++ {
		ret += string(chars[rand.Intn(len(chars))])
	}

	return ret
}

// Generate a random key (qBit 4.6.3)
func RandomKey() string {
	chars := "0123456789ABCDEF"
	ret := ""

	for i := 0; i < 8; i++ {
		ret += string(chars[rand.Intn(len(chars))])
	}

	return ret
}

func RandomPort() int {
	// TODO: maybe a fixed port on a specific machine?
	return rand.Intn(65535-1024) + 1024
}

func ParseAndRegenerateTorrent(filename string, fakeTracker string) (string, string, int64, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return "", "", 0, err
	}

	var m map[string]interface{}
	err = bencode.Unmarshal(fileContent, &m)
	if err != nil {
		return "", "", 0, err
	}

	// Calculate info_hash from torrent file
	infoValue, ok := m["info"]
	if !ok {
		return "", "", 0, errors.New("info field missing from torrent file")
	}

	infoBytes, err := bencode.Marshal(infoValue)
	if err != nil {
		return "", "", 0, err
	}

	hasher := sha1.New()
	hasher.Write(infoBytes)
	infoHash := hex.EncodeToString(hasher.Sum(nil))

	// Calculate total size of files
	infoMap, ok := infoValue.(map[string]interface{})
	if !ok {
		return "", "", 0, errors.New("info field is not a map")
	}

	var totalSize int64
	files, ok := infoMap["files"].([]interface{})
	if !ok {
		length, ok := infoMap["length"].(int64)
		if !ok {
			return "", "", 0, errors.New("files field and length field are missing")
		}

		totalSize = length
	} else {
		for _, file := range files {
			fileMap, ok := file.(map[string]interface{})
			if !ok {
				return "", "", 0, errors.New("file is not a map")
			}

			length, ok := fileMap["length"].(int64)
			if !ok {
				return "", "", 0, errors.New("length field is missing or is not an integer")
			}

			totalSize += length
		}
	}

	// Regenerate a new torrent with fake tracker
	origAnnounce := m["announce"].(string)
	m["announce"] = fakeTracker
	newTorrent, err := bencode.Marshal(m)
	if err != nil {
		return "", "", 0, err
	}
	err = os.WriteFile("FREE_"+filename, newTorrent, 0644)
	if err != nil {
		return "", "", 0, err
	}
	os.Remove(filename)

	return origAnnounce, infoHash, totalSize, nil
}
