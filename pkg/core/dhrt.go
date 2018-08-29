package dhrt

import (
	"github.com/satori/go.uuid"
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
	Key     string
	Value   string
	Version int64
	Created int64
	Updated int64
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
func (s *Store) Set(key string, value string) (Datum, error) {
	if key == "" {
		u := uuid.NewV4()
		key = u.String()
	}
	s.Lock()
	defer s.Unlock()
	now := time.Now().UnixNano()
	data := s.data
	datum := Datum{
		Value:   value,
		Updated: now,
	}
	// Check if the item exists, if not add a created timestamp
	_, exists := data[key]
	if !exists {
		datum.Created = now
		datum.Key = key
	}
	data[key] = datum
	s.data = data
	// Set here
	return datum, nil
}

// Del removes one key/value datum
func (s *Store) Del(key string) error {
	s.Lock()
	defer s.Unlock()
	delete(s.data, key)
	// TODO: remove item from cluster
	return nil
}

// DelAll removes an array of key/values
func (s *Store) DelAll(keys []string) error {
	for _, key := range keys {
		s.Del(key)
	}
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
