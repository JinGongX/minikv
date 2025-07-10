package main

import (
	"flag"
	"minikv/router"
	"minikv/storage"
)

func main() {
	mode := flag.String("mode", "router", "start mode: router or node")
	port := flag.String("port", "8080", "service port")
	flag.Parse()

	if *mode == "router" {
		router.StartRouter(*port, "./config/nodes.json")
	} else {
		store := storage.NewStore()
		store.Start(*port)
	}
}
