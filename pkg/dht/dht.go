package dht

import (
	"github.com/stevelacy/dhrt/pkg/core"
)

func XorDistance(a []byte, b []byte) []byte {
	var xor []byte
	xor = []byte{len(key)}
	for i, v := range a {
		xor = v ^ b[i]
	}
	return xor

}

func PeerDistance(n dhrt.Node, key []byte) []byte {
	return XorDistance(n.Id, key)
}
