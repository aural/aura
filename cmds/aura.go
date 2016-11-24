package main

import (
	"flag"
	"log"
	"time"

	zmq "github.com/pebbe/zmq4"
)

const (
	defaultTrackName = "data/example.flac"
	serverLocation   = "tcp://127.0.0.1:9090"
)

func init() {
	flag.Parse()
}

func main() {
	var trackLocations []string
	socket, err := zmq.NewSocket(zmq.REQ)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connecting to server at", serverLocation)
	socket.Connect(serverLocation)
	defer socket.Close()

	poller := zmq.NewPoller()
	poller.Add(socket, zmq.POLLIN)

	if flag.NArg() > 0 {
		trackLocations = flag.Args()
	} else {
		trackLocations = append(trackLocations, defaultTrackName)
	}

	for _, trackLocation := range trackLocations {
		log.Println("Requesting stream for", trackLocation)
		socket.SendMessage(trackLocation)

		sockets, err := poller.Poll(5000 * time.Millisecond)

		if err != nil {
			log.Println(err)
			break
		}

		if len(sockets) > 0 {
			reply, err := socket.RecvMessage(0)

			if err != nil {
				break
			}

			log.Println("-->", reply[0])
		}
	}
}
