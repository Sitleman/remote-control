package internal

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"remote-control/internal/models"
	"remote-control/internal/wscontrol"
	"remote-control/internal/wsweb"
	"time"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, it's remote-control server."))
}

func Index(w http.ResponseWriter, r *http.Request) {
	//http.Redirect(w, r, "/static/", http.StatusTemporaryRedirect)
	http.ServeFile(w, r, "static/index.html")
}

func InitRouter(logger *zap.Logger, db *badger.DB) *http.Server {
	_, webUpdates := make(chan models.Update), make(chan models.Update)
	wsc := wscontrol.NewWSControl(logger, db, webUpdates)
	wsw := wsweb.NewWSWeb(logger, db, wsc.Updates, webUpdates)

	r := mux.NewRouter()
	r.HandleFunc("/", Index)
	r.HandleFunc("/me", MeHandler)
	r.HandleFunc("/ws-control", wsc.Handler)
	r.HandleFunc("/ws-web", wsw.Handler)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	return &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 5 * time.Minute,
		ReadTimeout:  5 * time.Minute,
	}
}
