package services

type Hub struct {
	// todo rwmutex
	// connected clients; room is nil if a client didn't join any
	clients map[*Client]*Room
	// a map of taken nicknames
	nicknames map[string]struct{}
	// available rooms to join
	rooms map[string]*Room
}

func NewHub() *Hub {
	return &Hub{
		clients:   make(map[*Client]*Room),
		nicknames: make(map[string]struct{}),
		rooms: map[string]*Room{
			"general": {clients: make(map[*Client]struct{})},
			"random":  {clients: make(map[*Client]struct{})},
		},
	}
}

func (s *Hub) Add(conn WS) {
	client := NewClient(conn, s, conn.RemoteAddr().String())

	s.clients[client] = nil

	go client.ReadMessages()
}

func (s *Hub) Remove(client *Client) {
	room := s.clients[client]
	if room != nil {
		room.Leave(client)
	}

	delete(s.clients, client)
	delete(s.nicknames, client.nickname)
}

// todo make func to return (*ResponseMessage, error) and make the caller deal w/ the response
func (s *Hub) ProcessClientMessage(client *Client, message RequestMessage) {
	// todo better routing than switch-case. setup funcs in map etc
	switch message.Action {
	case "SET_NICKNAME":
		ActionSetNickname(s, client, message)

	case "LIST_ROOMS":
		ActionListRooms(s, client, message)

	case "SET_ROOM":
		ActionSetRoom(s, client, message)

	case "CHAT":
		ActionChat(s, client, message)

	case "GET_USERS_IN_ROOM":
		ActionGetUsersInRoom(s, client, message)

	default:
		ActionDefault(s, client, message)
	}
}

func ActionDefault(_ *Hub, client *Client, _ RequestMessage) {
	client.SendMessage(ResponseMessage{Error: &Error{Code: "UNKNOWN_ACTION"}})
}
