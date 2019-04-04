package services

func ActionChat(hub *Hub, client *Client, message RequestMessage) {
	if len(client.nickname) == 0 {
		client.SendMessage(ResponseMessage{Error: &Error{Code: "NO_NICKNAME"}})
		return
	}

	room := hub.clients[client]
	if room == nil {
		client.SendMessage(ResponseMessage{Error: &Error{Code: "NOT_JOINED_TO_ROOM"}})
		return
	}

	room.BroadcastMessage(ResponseMessage{Type: "CHAT", Message: client.nickname + ": " + message.Message})
}
