package dhrt

import (
	"bytes"
	"os"
	"sync"
	"time"
)

// Config for the db initialization
type Config struct {
	Root  string
	Nodes []Node
	Node  Node
}

// Node is a single DB in a cluster
type Node struct {
	ListenAddr   string // The address to listen on
	AnnounceAddr string // The address for the other nodes to make requests
	Id           []byte // The id of the node
	Port         int
	Version      int
	Status       string // healthy, unhealthy, dead
}

// Store core
type Store struct {
	data   map[string]Datum
	Config Config
	sync.RWMutex
	file     *os.File // TODO
	Open     bool
	CaughtUp bool // TODO
}

// Stats core
type Stats struct {
	Config    Config
	Size      int // Size of the data
	ItemCount int // Number of items in the store
}

// Datum is the individual item
type Datum struct {
	Key     []byte
	Value   string
	Version int64
	Created int64
	Updated int64
	Hash    []byte
}

// Open the database
func Open(config Config) (*Store, error) {
	store := &Store{
		data:   map[string]Datum{},
		Config: config,
		Open:   true,
	}
	return store, nil
}

// Close clears the db of all memory items
func (s *Store) Close() error {
	s.data, s.file = nil, nil
	s.Open = false
	return nil
}

// Get a single value with a key
// This method will lock the mutex
func (s *Store) Get(key string) (Datum, error) {
	s.RLock()
	defer s.RUnlock()
	datum := s.data[key]
	return datum, nil
}

// Set a single value with a key
// Returns an id and error
// This method will lock the mutex
func (s *Store) Set(key []byte, value string) (Datum, error) {
	var currentHash []byte = Hash(value) // always hash the data
	if key == nil || bytes.Equal(key, []byte{}) {
		key = currentHash
	}
	s.Lock()
	defer s.Unlock()
	now := time.Now().UnixNano()
	data := s.data
	datum := Datum{
		Value:   value,
		Updated: now,
		Hash:    currentHash,
	}
	// Check if the item exists, if not add a created timestamp
	dKey := HashToString(key)
	_, exists := data[dKey]
	if !exists {
		datum.Created = now
		datum.Key = key
	}
	data[dKey] = datum
	s.data = data // Set here
	return datum, nil
}

// Del removes one key/value datum
func (s *Store) Del(key []byte) error {
	s.Lock()
	defer s.Unlock()
	delete(s.data, HashToString(key))
	return nil
}

// Stats returns the status of the node's store
func (s *Store) Stats() Stats {
	stats := Stats{
		ItemCount: len(s.data),
		Config:    s.Config,
	}
	return stats
}
