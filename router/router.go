package router

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"minikv/hash"
)

var ring *hash.HashRing

func LoadNodes(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var nodes []string
	err = json.Unmarshal(data, &nodes)
	return nodes, err
}

func StartRouter(port string, configPath string) {
	nodes, err := LoadNodes(configPath)
	if err != nil {
		log.Fatal("Error loading nodes:", err)
	}

	ring = hash.New(3)
	ring.Add(nodes...)

	http.HandleFunc("/put", handlePut)
	http.HandleFunc("/get", handleGet)

	fmt.Println("Router started at port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func forward(method, node, key, value string) (string, error) {
	var url string
	if method == "PUT" {
		url = fmt.Sprintf("%s/put?key=%s&value=%s", node, key, value)
	} else {
		url = fmt.Sprintf("%s/get?key=%s", node, key)
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	res, err := io.ReadAll(resp.Body)
	return string(res), err
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	node := ring.Get(key)
	resp, err := forward("PUT", node, key, value)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte(resp))
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	node := ring.Get(key)
	resp, err := forward("GET", node, key, "")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte(resp))
}
