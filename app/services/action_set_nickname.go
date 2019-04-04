package services

func ActionSetNickname(hub *Hub, client *Client, message RequestMessage) {
	if len(message.Message) == 0 {
		client.SendMessage(ResponseMessage{Error: &Error{Code: "NO_NICKNAME"}})
		return
	}

	if _, exists := hub.nicknames[message.Message]; exists {
		client.SendMessage(ResponseMessage{Error: &Error{Code: "NICKNAME_TAKEN"}})
		return
	}

	if len(client.nickname) > 0 {
		delete(hub.nicknames, client.nickname)
		client.nickname = ""
	}

	client.nickname = message.Message
	hub.nicknames[client.nickname] = struct{}{}
	client.SendMessage(ResponseMessage{Message: "nickname set"})
}
