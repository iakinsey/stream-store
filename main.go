package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/iakinsey/stream-store/controller"
	"github.com/iakinsey/stream-store/util"
)

func main() {
	addr := util.Environ("STREAM_STORE_ADDRESS", "0.0.0.0:40865")
	sm := http.NewServeMux()

	sm.HandleFunc("/", controller.Router)

	server := &http.Server{
		Addr:    addr,
		Handler: sm,
	}
	fmt.Printf("Stream server listening on %s\n", addr)
	log.Fatal(server.ListenAndServe())
}
