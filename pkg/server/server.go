package server

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/stevelacy/dhrt/pkg/core"
	"net/http"
)

// Start implements the http service
func Start(db *dhrt.Store, config dhrt.Config) error {
	router := httprouter.New()
	router.GET("/", Index(db))
	router.GET("/:key", GetItem(db))
	router.POST("/:key/:value", SetItem(db))
	router.DELETE("/:key", DelItem(db))
	fmt.Printf("Server starting on %s port %d\n", config.Node.ListenAddr, config.Node.Port)
	addr := fmt.Sprintf("%s:%d", config.Node.ListenAddr, config.Node.Port)
	return http.ListenAndServe(addr, router)
}

// Index lists stats, should also list collections
func Index(db *dhrt.Store) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		data, err := json.Marshal(db.Stats())
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"Error": "%e"}`, err)))
			w.WriteHeader(500)
			return
		}
		w.Write(data)
	}
}

// GetItem returns one key/value
func GetItem(db *dhrt.Store) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		res, err := db.Get(p.ByName("key"))
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"Error": "%e"}`, err)))
			w.WriteHeader(500)
			return
		}
		data, err := json.Marshal(res)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"Error": "%e"}`, err)))
			w.WriteHeader(500)
			return
		}
		w.Write(data)
	}
}

// SetItem sets the value of one key
func SetItem(db *dhrt.Store) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		res, err := db.Set(p.ByName("key"), p.ByName("value"))
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"Error": "%e"}`, err)))
			w.WriteHeader(500)
			return
		}
		data, err := json.Marshal(res)
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"Error": "%e"}`, err)))
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(201)
		w.Write(data)
	}
}

// DelItem deletes one key/value
func DelItem(db *dhrt.Store) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		key := p.ByName("key")
		err := db.Del(key)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.Write([]byte(fmt.Sprintf(`{"Error": "%e"}`, err)))
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(202)
		w.Write([]byte(fmt.Sprintf(`{"Key": "%s"}`, key)))
	}
}
