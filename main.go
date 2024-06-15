package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connPool *pgx.Conn

func initDB() {
	var err error
	connPool, err = pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/real_time_db")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	log.Println("Connected to the database successfully!")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	go func() {
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Println("ReadMessage error:", err)
				return
			}
			if err := conn.WriteMessage(messageType, p); err != nil {
				log.Println("WriteMessage error:", err)
				return
			}
		}
	}()

	log.Println("Listening for database notifications...")
	if _, err := connPool.Exec(context.Background(), "LISTEN events"); err != nil {
		log.Fatalf("Failed to listen for notifications: %v", err)
	}

	for {
		notification, err := connPool.WaitForNotification(context.Background())
		if err != nil {
			log.Println("WaitForNotification error:", err)
			continue
		}
		log.Printf("Received notification: %v", notification.Payload)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(notification.Payload)); err != nil {
			log.Println("WriteMessage error:", err)
			return
		}
	}
}

func main() {
	initDB()
	defer connPool.Close(context.Background())

	http.HandleFunc("/ws", wsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
