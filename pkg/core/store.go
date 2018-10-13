package dhrt

type StoreIf interface {
	Contains(key []byte) bool                    // _R__
	Get(key []byte) (Datum, error)               // _R__
	Set(key []byte, value string) (Datum, error) // C_U_
	Del(Key []byte) error                        // ___D
}

type LocalStoreIf interface {
	StoreIf
	DelAll(keys [][]byte) error
	Stats() Stats
}
