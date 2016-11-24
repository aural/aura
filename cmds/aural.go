package main

import (
	"flag"
	"log"

	"github.com/aural/aural"

	zmq "github.com/pebbe/zmq4"
)

const (
	serverLocation = "tcp://127.0.0.1:9090"
)

func init() {
	flag.Parse()
}

func tracks(locations []string) (tracks []*aural.Track) {
	var track *aural.Track

	for _, location := range locations {
		track = new(aural.Track)
		track.Location = location
		tracks = append(tracks, track)
	}

	return tracks
}

func createServer(playstate *aural.Playstate) chan string {
	channel := make(chan string)
	server, err := zmq.NewSocket(zmq.REP)

	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		defer server.Close()

		server.Bind(serverLocation)
		log.Println("Listening at", serverLocation)

		for {
			request, err := server.RecvMessage(0)

			if err != nil {
				log.Fatalln(err)
			}

			log.Println("<-", request[0])
			playstate.Playlist.Queue(request[0])
			server.SendMessage("OK")

			channel <- request[0]
		}
	}()

	return channel
}

func main() {
	log.Println("Starting aural daemon")
	defer aural.Terminate()

	playstate := aural.NewPlaystate()

	go createServer(playstate)
	audio := playstate.MainLoop()

	for {
		<-audio
	}
}
