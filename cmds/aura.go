package main

import (
	"flag"
	"log"
	"time"

	"github.com/aural/aural"
	zmq "github.com/pebbe/zmq4"
)

const (
	defaultTrackName = "data/example.flac"
)

type HandlerMap map[string]func(*zmq.Socket, []string)

var handlers HandlerMap

func init() {
	flag.Parse()

	handlers = make(HandlerMap)
	handlers["load"] = LoadHandler
	handlers["clear"] = ClearHandler
}

func LoadHandler(socket *zmq.Socket, trackIdentifiers []string) {
	poller := zmq.NewPoller()
	poller.Add(socket, zmq.POLLIN)

	if len(trackIdentifiers) == 0 {
		trackIdentifiers = append(trackIdentifiers, defaultTrackName)
	}

	for _, trackIdentifier := range trackIdentifiers {
		socket.SendMessage("LOAD", trackIdentifier)
		sockets, err := poller.Poll(5000 * time.Millisecond)

		if err != nil {
			log.Println(err)
			break
		}

		if len(sockets) > 0 {
			_, err := socket.RecvMessage(0)

			if err != nil {
				log.Println("Error with ", trackIdentifier, ": ", err)
			}
		}
	}
}

func ClearHandler(socket *zmq.Socket, trackIdentifiers []string) {
	poller := zmq.NewPoller()
	poller.Add(socket, zmq.POLLIN)

	socket.SendMessage("CLEAR")
	sockets, err := poller.Poll(5000 * time.Millisecond)

	if err != nil {
		log.Println(err)
		return
	}

	if len(sockets) > 0 {
		_, err := socket.RecvMessage(0)

		if err != nil {
			log.Println("Error clearing playstate. ", err)
		}
	}
}

func main() {
	configuration := aural.GetConfiguration()
	socket, err := zmq.NewSocket(zmq.REQ)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connecting to server at", configuration.Address)
	socket.Connect(configuration.Address)
	defer socket.Close()

	if flag.NArg() == 0 {
		log.Fatalln("Play/Pause toggle not yet implemented")
	}

	arguments := flag.Args()
	handler, ok := handlers[arguments[0]]

	if ok {
		arguments = arguments[1:]
	} else {
		handler, ok = handlers["load"]

		if !ok {
			log.Fatalln("Could not parse command-line. Sorry! :'('")
		}
	}

	handler(socket, arguments)
}
