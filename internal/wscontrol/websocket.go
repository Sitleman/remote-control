package wscontrol

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"net/http"
	"remote-control/internal/models"
)

type WSControl struct {
	connects     map[string]*websocket.Conn
	db           *badger.DB
	logger       *zap.Logger
	inMessages   chan Message
	sessionCount int

	Updates   chan string
	inUpdates chan models.Update
}

type Message struct {
	Type int
	Text string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewWSControl(logger *zap.Logger, db *badger.DB, inUpdates chan models.Update) *WSControl {
	wsServer := &WSControl{
		connects:   map[string]*websocket.Conn{},
		db:         db,
		logger:     logger,
		inMessages: make(chan Message),
		Updates:    make(chan string),
		inUpdates:  inUpdates,
	}
	go wsServer.run()
	return wsServer
}

func (wsc *WSControl) Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("wscontrol request")

	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
		r.Body.Close()
	}()

	//uuid := r.Header.Get("UUID")
	//if uuid == "" {
	//	wsc.logger.Warn("wscontrol request must have UUID")
	//	w.WriteHeader(http.StatusBadRequest)
	//	w.Write([]byte("wscontrol request must have UUID"))
	//	return
	//}

	con, _ := upgrader.Upgrade(w, r, nil)
	wsc.connects[fmt.Sprintf("%d", wsc.sessionCount)] = con
	wsc.sessionCount++

	go func() {
		for {
			mt, message, err := con.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			wsc.logger.Info("New message", zap.Int("Type", mt),
				zap.String("Message", string(message)))
			wsc.inMessages <- Message{mt, string(message)}
		}
	}()

	wsc.logger.Info("New ws-control connection")
	//session, err := r.Cookie("session")
	//if err != nil {
	//	log.Printf("fail to get session: %v", err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//}
	//
}

func (wsc *WSControl) run() {
	//var message Message
	for {
		select {
		case _ = <-wsc.inMessages:
			newColor := fmt.Sprintf("#%02x%02x%02x", rand.Intn(255), rand.Intn(255), rand.Intn(255))
			wsc.Updates <- newColor
		case update := <-wsc.inUpdates:
			wsc.updateAllClients(1, update)
			wsc.logger.Info("Send update control")
		}

	}

}

func (wsc *WSControl) updateAllClients(mt int, update models.Update) {
	updateJson, err := json.Marshal(update)
	if err != nil {
		wsc.logger.Warn("Fail marshal update", zap.Error(err))
		return
	}
	for session, conn := range wsc.connects {
		err := conn.WriteMessage(mt, updateJson)
		if err != nil {
			wsc.logger.Warn("Fail to send update message for device", zap.String("session", session),
				zap.Error(err))
			continue
		}
		wsc.logger.Warn("Successful send update for device", zap.String("session", session))
	}
}
