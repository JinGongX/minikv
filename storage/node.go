package storage

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Store struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

func (s *Store) Start(port string) {
	http.HandleFunc("/put", s.handlePut)
	http.HandleFunc("/get", s.handleGet)

	fmt.Println("Storage node started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (s *Store) handlePut(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	s.mu.Lock()
	s.data[key] = value
	s.mu.Unlock()
	w.Write([]byte("OK"))
}

func (s *Store) handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	s.mu.RLock()
	val, ok := s.data[key]
	s.mu.RUnlock()
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte(val))
}
