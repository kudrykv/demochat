package services

import "strings"

func ActionGetUsersInRoom(hub *Hub, client *Client, _ RequestMessage) {
	room := hub.clients[client]
	if room == nil {
		client.SendMessage(ResponseMessage{Error: &Error{Code: "NOT_JOINED_TO_ROOM"}})
		return
	}

	nicknames := room.GetClientsNicknames()
	client.SendMessage(ResponseMessage{Message: "users in room: " + strings.Join(nicknames, ", ")})
}
