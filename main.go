package main

import (
	"log"
	"net"

	"github.com/iakinsey/stream-store/protocol"
	"github.com/iakinsey/stream-store/util"
)

func main() {
	addr := util.Environ("STREAM_STORE_ADDRESS", "0.0.0.0:40865")

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("Stream store listening on %s\n", addr)
	log.Fatalf(protocol.Listen(listener).Error())
}
