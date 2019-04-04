package services

import "strings"

func ActionListRooms(hub *Hub, client *Client, _ RequestMessage) {
	if len(client.nickname) == 0 {
		client.SendMessage(ResponseMessage{Error: &Error{Code: "NO_NICKNAME"}})
		return
	}

	rooms := make([]string, 0, len(hub.rooms))
	for k := range hub.rooms {
		rooms = append(rooms, k)
	}

	client.SendMessage(ResponseMessage{Message: "available rooms: " + strings.Join(rooms, ", ")})
}
