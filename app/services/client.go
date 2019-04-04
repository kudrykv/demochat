package services

import (
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	conn    WS
	headers http.Header
	ip      string

	nickname string

	hub *Hub
}

type WS interface {
	ReadJSON(v interface{}) error
	WriteJSON(v interface{}) error
	RemoteAddr() net.Addr
	Close() error
}

// Creates new client
// IP is solely for logging
func NewClient(conn WS, hub *Hub, ip string) *Client {
	return &Client{
		conn: conn,
		hub:  hub,
		ip:   ip,
	}
}

// ReadMessages listens for messages from the client
func (c *Client) ReadMessages() {
	defer func(c *Client) {
		c.hub.Remove(c)
		_ = c.conn.Close()

		log.WithFields(log.Fields{
			"nickname": c.nickname,
			"ip":       c.ip,
		}).Info("client disconnected")
	}(c)

	idleTime := 2 * time.Minute
	idleTimer := time.NewTimer(idleTime)
	msgChan := make(chan RequestMessage) // messages from user
	exit := make(chan bool)              // read-message loop breaker

	go readMsgOrExit(c, msgChan, exit)

	for {
		l := log.WithFields(log.Fields{
			"nickname": c.nickname,
			"ip":       c.ip,
		})

		select {
		case <-idleTimer.C:
			l.Info("user idle, disconnect")
			return

		case msg := <-msgChan:
			idleTimer.Reset(idleTime)

			l = l.WithField("msg", msg)

			l.Info("processing user message")
			c.hub.ProcessClientMessage(c, msg)
			l.Info("processed user message")
		}
	}

	// todo implement ping-pong, set timeouts, limits etc
}

func readMsgOrExit(u *Client, msgChan chan RequestMessage, exit chan bool) {
	for {
		select {
		case <-exit:
			log.WithFields(log.Fields{
				"nickname": u.nickname,
				"ip":       u.ip,
			}).Info("exiting message read loop")
			return

		default:
			var msg RequestMessage
			if err := u.conn.ReadJSON(&msg); err != nil {
				log.WithFields(log.Fields{
					"nickname": u.nickname,
					"ip":       u.ip,
					"err":      err,
				}).Warn("not JSON or closed")
				return
			}

			msgChan <- msg
		}
	}
}

// SendMessage to the user
func (c *Client) SendMessage(resp ResponseMessage) {
	resp.Timestamp = time.Now().Round(time.Millisecond)

	l := log.WithFields(log.Fields{
		"nickname": c.nickname,
		"ip":       c.ip,
		"resp":     resp,
	})

	if resp.Error != nil {
		l = l.WithField("resp_err", resp.Error)
	}

	if err := c.conn.WriteJSON(resp); err != nil {
		l.WithField("err", err).Error("failed to write message")
		c.hub.Remove(c)
	} else {
		l.Info("wrote msg")
	}
}

func (c *Client) Nickname() string {
	return c.nickname
}

func (c *Client) SetNickname(nickname string) {
	c.nickname = nickname
}
