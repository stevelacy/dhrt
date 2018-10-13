package dhrt

import "errors"

func FindPeer(key []byte) StoreIf {
	return nil
}

// DelAll removes an array of key/values
func (s *Store) DelAll(keys [][]byte) error {
	for _, key := range keys {
		if s.Contains(key) {
			s.Del(key)
		} else {
			var peer StoreIf
			peer = FindPeer(key)
			if peer == nil {
				return errors.New("No peer for key")
			}
			peer.Del(key)
		}
	}
	return nil
}
