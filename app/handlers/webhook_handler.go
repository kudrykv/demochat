package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kudrykv/demochat/app/services"
	log "github.com/sirupsen/logrus"
)

type WebsocketHandler struct {
	ws websocket.Upgrader

	hub *services.Hub
}

func NewWebsocketHandler(hub *services.Hub) *WebsocketHandler {
	return &WebsocketHandler{
		hub: hub,
		ws:  websocket.Upgrader{},
	}
}

func (h *WebsocketHandler) Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := h.ws.Upgrade(w, r, nil)
	if err != nil {
		log.WithField("err", err).Error("failed to upgrade the connection")
		return
	}

	h.hub.Add(conn)

	log.WithFields(log.Fields{
		"user-agent": r.Header.Get("User-Agent"),
		"ip":         r.RemoteAddr,
	}).Info("connected client")
}
