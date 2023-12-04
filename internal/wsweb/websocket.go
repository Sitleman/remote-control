package wsweb

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"
	"net/http"
	"remote-control/internal/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type webMessage struct {
	Type    int
	Element string `json:"type"`
	Data    string `json:"data"`
}

type WSWeb struct {
	connects     map[string]*websocket.Conn
	db           *badger.DB
	logger       *zap.Logger
	sessionCount int

	inMessages chan webMessage
	outUpdates chan models.Update
	inUpdates  chan string
}

func NewWSWeb(logger *zap.Logger, db *badger.DB, inUpdates chan string, outUpdate chan models.Update) *WSWeb {
	server := &WSWeb{
		connects:   map[string]*websocket.Conn{},
		db:         db,
		logger:     logger,
		inMessages: make(chan webMessage),
		outUpdates: outUpdate,
		inUpdates:  inUpdates,
	}
	go server.run()
	return server
}

func (wsw *WSWeb) Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("wscontrol request")

	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
		r.Body.Close()
	}()

	//session, err := r.Cookie("session")
	//if err != nil {
	//	log.Printf("fail to get session: %v", err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//}

	con, _ := upgrader.Upgrade(w, r, nil)
	wsw.connects[fmt.Sprintf("%d", wsw.sessionCount)] = con
	wsw.sessionCount++

	go func() {
		for {
			mt, rawMess, err := con.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			var message webMessage
			err = json.Unmarshal(rawMess, &message)
			if err != nil {
				wsw.logger.Warn("Incorrect new message", zap.Int("Type", mt),
					zap.String("Message", string(rawMess)), zap.Error(err))
				continue
			}

			wsw.logger.Info("New message", zap.Int("Type", mt), zap.Any("Message", message))
			wsw.inMessages <- message
		}
	}()

	wsw.logger.Info("New ws-web connection")
}

func (wsw *WSWeb) run() {
	for {
		select {
		case newColor := <-wsw.inUpdates:
			wsw.updateAllClients(1, newColor)
		case message := <-wsw.inMessages:
			wsw.outUpdates <- models.Update{
				Device:  "",
				Element: message.Element,
				Data:    message.Data,
			}
			wsw.logger.Info("Send update web")
		}
	}
}

func (wsw *WSWeb) updateAllClients(mt int, message string) {
	for session, conn := range wsw.connects {
		err := conn.WriteMessage(mt, []byte(message))
		if err != nil {
			wsw.logger.Warn("Fail to send update message for client", zap.String("session", session),
				zap.Error(err))
			continue
		}
		wsw.logger.Warn("Successful send update for client", zap.String("session", session))
	}
}

//func (wsc *WSControl) run() {
//
//		log.Printf("recv: %s", message)
//		err = c.WriteMessage(mt, message)
//		if err != nil {
//			log.Println("write:", err)
//			break
//		}
//	}
//}
//
//func SendMessage(session string, message string) error {
//	con, ok := savedsockets[session]
//	if !ok {
//		log.Printf("websocket: no such session in ws storage")
//		return nil
//	}
//	if err := con.WriteMessage(1, []byte(message)); err != nil {
//		log.Printf("websocket: fail to write message: %v", err)
//	}
//	return nil
//}
