package main

import (
	"net/http"
	"os"

	"github.com/kudrykv/demochat/app/handlers"
	"github.com/kudrykv/demochat/app/services"
	log "github.com/sirupsen/logrus"
	"goji.io"
	"goji.io/pat"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{ForceColors: true})

	mux := goji.NewMux()

	hubSvc := services.NewHub()
	whHandler := handlers.NewWebsocketHandler(hubSvc)

	mux.HandleFunc(pat.New("/ws"), whHandler.Websocket)
	mux.HandleFunc(pat.Get("/"), func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./app/static/home.html")
	})

	log.Info("server starting, http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
