package server

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type WebSocketManager struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

var wsManager = WebSocketManager{
	clients: make(map[*websocket.Conn]bool),
}

func Run(api fiber.Router) error {
	api.Get("/message", getAllMessages)
	api.Use("/ws", initWebsockets)
	api.Get("/ws", websocket.New(connectWS))
	return nil
}

func getAllMessages(c *fiber.Ctx) error {
	messgeCollection, err := GetMessages()
	if err != nil {
		log.Println("Error fetching messages: ", err)
	}
	return c.Status(200).JSON(messgeCollection)
}

func initWebsockets(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func connectWS(c *websocket.Conn) {
	wsManager.mu.Lock()
	wsManager.clients[c] = true
	wsManager.mu.Unlock()

	defer func() {
		wsManager.mu.Lock()
		delete(wsManager.clients, c)
		wsManager.mu.Unlock()
		c.Close()
	}()

	err := sendSnapshot(c)
	if err != nil {
		return
	}

	for {
		err = handle(c)
		if err != nil {
			break
		}
	}
}

func broadcast(message []byte) {
	wsManager.mu.Lock()
	defer wsManager.mu.Unlock()

	for client := range wsManager.clients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			client.Close()
			delete(wsManager.clients, client)
		}
	}
}

func sendSnapshot(c *websocket.Conn) error {
	messages, err := GetMessages()
	if err != nil {
		fmt.Println("Failed to get all messages for snapshot!")
		return err
	}

	snapshotJson, _ := json.Marshal(messages)

	if err := c.WriteMessage(websocket.TextMessage, snapshotJson); err != nil {
		return err
	}

	return nil
}

func handle(c *websocket.Conn) error {
	_, msg, err := c.ReadMessage()
	if err != nil {
		return err
	}

	var message Message
	if err := json.Unmarshal(msg, &message); err != nil {
		log.Println("Error parsing message: ", err)
		c.WriteMessage(websocket.TextMessage, []byte("Message handling failed"))
		return nil
	}

	if message.Content == "" {
		log.Println("Received empty message, ignoring")
		c.WriteMessage(websocket.TextMessage, []byte("Message can not be empty string"))
		return nil
	}

	message, err = AddMessage(message)
	if err != nil {
		return err
	}

	newMsgJson, err := json.Marshal(message)
	if err != nil {
		return err
	}
	broadcast(newMsgJson)

	return nil
}
