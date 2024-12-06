package server

type Messenger struct {
	broadcast chan []byte
}

func newMessenger() *Messenger {
	return &Messenger{
		broadcast: make(chan []byte),
	}
}
