package dht

import (
	"bytes"
)

const LENGTH = 4

type Key struct {
	Contents []int64
}

func (k *Key) Equal(b Key) bool {
	return bytes.Equal(k.Contents, b.Contents)
}
