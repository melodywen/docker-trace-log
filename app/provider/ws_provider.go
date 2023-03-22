package provider

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/melodywen/docker-trace-log/contracts"
	"log"
	"net/http"
)

type WsProvider struct {
}

func NewWsProvider() *WsProvider {
	return &WsProvider{}
}

func (w WsProvider) StartServerBeforeEvent(ctx context.Context, app contracts.AppAttributeInterface) error {
	defer app.GetLog().EnterExitFunc(ctx)()

	return nil
}

var upgrader = websocket.Upgrader{} // use default options

func (w WsProvider) StartServerAfterEvent(ctx context.Context, app contracts.AppAttributeInterface) error {
	defer app.GetLog().EnterExitFunc(ctx)()

	/**
	  增加ws
	*/
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		// Upgrade our raw HTTP connection to a websocket based one
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Print("Error during connection upgradation:", err)
			return
		}
		defer conn.Close()

		// The event loop
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error during message reading:", err)
				break
			}
			log.Printf("Received: %s", message)
			err = conn.WriteMessage(messageType, message)
			if err != nil {
				log.Println("Error during message writing:", err)
				break
			}
		}
	})

	return nil
}
