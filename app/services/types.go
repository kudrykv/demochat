package services

import "time"

type RequestMessage struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type ResponseMessage struct {
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Error     *Error    `json:"error"`
	Timestamp time.Time `json:"timestamp"`
}

type Error struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}
