package aural

type MessageHandler func(*Playstate, []string)
type handlerMap map[string]MessageHandler

const (
	DEFAULT_PORT int = 28346
)

var handlers handlerMap

func init() {
	handlers = make(handlerMap)

	handlers["LOAD"] = LoadHandler
	handlers["CLEAR"] = ClearHandler
}

func LoadHandler(playstate *Playstate, arguments []string) {
	playstate.Playlist.Queue(arguments[0])
}

func ClearHandler(playstate *Playstate, arguments []string) {
	playstate.Clear()
}

func HandleRequest(playstate *Playstate, kind string, arguments []string) {
	handler, ok := handlers[kind]

	if !ok {
		return
	}

	handler(playstate, arguments)
}
