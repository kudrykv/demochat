package services

type Room struct {
	// todo rwmutex
	// list of clients connected to the room
	clients map[*Client]struct{}
}

func (r *Room) BroadcastMessage(message ResponseMessage) {
	for k := range r.clients {
		k.SendMessage(message)
	}
}

func (r *Room) NotifyJoin(nickname string) {
	for k := range r.clients {
		k.SendMessage(ResponseMessage{Type: "BROADCAST", Message: nickname + " joined the room!"})
	}
}

func (r *Room) Leave(client *Client) {
	delete(r.clients, client)
	r.BroadcastMessage(ResponseMessage{Type: "BROADCAST", Message: "left the room"})
}

func (r *Room) GetClientsNicknames() []string {
	nicknames := make([]string, 0, len(r.clients))

	for k := range r.clients {
		nicknames = append(nicknames, k.nickname)
	}

	return nicknames
}
