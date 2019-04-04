package services_test

import (
	"testing"

	"github.com/kudrykv/demochat/app/services"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

// todo other tests

func TestActionSetNickname(t *testing.T) {
	Convey("ActionSetNickname", t, func() {
		hub := services.NewHub()

		ws := &services.WSMock{}
		client := services.NewClient(ws, hub, "::1")

		Convey("when the user has no nickname", func() {
			Convey("when the user passes the message", func() {
				Convey("then it is taken as a nickname", func() {
					msg := services.RequestMessage{Message: "nick"}

					ws.On("WriteJSON", mock.MatchedBy(nickSet)).Return(nil).Once()

					services.ActionSetNickname(hub, client, msg)

					So(client.Nickname(), ShouldEqual, "nick")

					mock.AssertExpectationsForObjects(t, ws)
				})
			})

			Convey("when the user passes no message", func() {
				Convey("then the error response should happen", func() {
					msg := services.RequestMessage{}

					ws.On("WriteJSON", mock.MatchedBy(emptyNick)).Return(nil).Once()

					services.ActionSetNickname(hub, client, msg)

					mock.AssertExpectationsForObjects(t, ws)
				})
			})

			Convey("when the nickname already exists", func() {
				msg := services.RequestMessage{Message: "nick"}

				separateWSMock := &services.WSMock{}
				separateWSMock.On("WriteJSON", mock.MatchedBy(nickSet)).Return(nil).Once()
				separateClient := services.NewClient(separateWSMock, hub, "::2")

				services.ActionSetNickname(hub, separateClient, msg)

				mock.AssertExpectationsForObjects(t, separateWSMock)

				Convey("then the err response should happen", func() {
					ws.On("WriteJSON", mock.MatchedBy(nickTaken)).Return(nil).Once()

					services.ActionSetNickname(hub, client, msg)

					mock.AssertExpectationsForObjects(t, ws)
				})
			})
		})

		Convey("when the user has nickname", func() {
			Convey("when the user specifies other nickname", func() {
				Convey("then current should be freed", func() {
					ws.On("WriteJSON", mock.MatchedBy(nickSet)).Return(nil).Twice()

					services.ActionSetNickname(hub, client, services.RequestMessage{Message: "nick"})

					So(client.Nickname(), ShouldEqual, "nick")

					services.ActionSetNickname(hub, client, services.RequestMessage{Message: "dock"})

					So(client.Nickname(), ShouldEqual, "dock")

					separateWSMock := &services.WSMock{}
					separateWSMock.On("WriteJSON", mock.MatchedBy(nickSet)).Return(nil).Once()
					separateClient := services.NewClient(separateWSMock, hub, "::2")

					services.ActionSetNickname(hub, separateClient, services.RequestMessage{Message: "nick"})

					mock.AssertExpectationsForObjects(t, ws, separateWSMock)
				})
			})
		})
	})
}

func nickSet(r services.ResponseMessage) bool {
	return r.Message == "nickname set"
}

func emptyNick(r services.ResponseMessage) bool {
	return r.Error != nil && r.Error.Code == "NO_NICKNAME"
}

func nickTaken(r services.ResponseMessage) bool {
	return r.Error != nil && r.Error.Code == "NICKNAME_TAKEN"
}
