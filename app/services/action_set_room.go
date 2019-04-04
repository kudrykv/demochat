package services

func ActionSetRoom(hub *Hub, client *Client, message RequestMessage) {
	if len(client.nickname) == 0 {
		client.SendMessage(ResponseMessage{Error: &Error{Code: "NO_NICKNAME"}})
		return
	}

	if len(message.Message) == 0 {
		client.SendMessage(ResponseMessage{Error: &Error{Code: "ROOM_NAME_EMPTY"}})
		return
	}

	room := hub.clients[client]
	if room != nil {
		room.Leave(client)
		hub.clients[client] = nil
		client.SendMessage(ResponseMessage{Message: "left the previous room"})
	}

	room, exists := hub.rooms[message.Message]
	if !exists {
		client.SendMessage(ResponseMessage{Error: &Error{Code: "ROOM_DOES_NOT_EXIST"}})
		return
	}

	room.clients[client] = struct{}{}
	hub.clients[client] = room
	room.NotifyJoin(client.nickname)
	client.SendMessage(ResponseMessage{Message: "room set to " + message.Message})
}
