package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins (lock down in prod)
	},
}

type WSMessage struct {
	Event    string      `json:"event"`
	Message  string      `json:"message,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	ClientID string      `json:"clientID,omitempty"`
}

type Client struct {
	ID       string
	Conn     *websocket.Conn
	Incoming chan map[string]interface{}
	WriteMu  sync.Mutex
}

var (
	clients   = make(map[string]*Client)
	clientsMu sync.Mutex
)

func (c *Client) readLoop() {
	defer func() {
		close(c.Incoming)
		c.Conn.Close()

		clientsMu.Lock()
		delete(clients, c.ID)
		clientsMu.Unlock()

		log.Println("Client disconnected:", c.ID)
	}()

	for {
		var msg map[string]interface{}
		if err := c.Conn.ReadJSON(&msg); err != nil {
			return
		}
		c.Incoming <- msg
	}
}

func serveStatic() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
}

func helloGet(w http.ResponseWriter, r *http.Request) {
	log.Println("Headers:", r.Header)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Hello from Go!",
	})
}

func helloPost(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	json.NewDecoder(r.Body).Decode(&body)

	log.Println("Headers:", r.Header)
	log.Println("Body:", body)

	resp := map[string]interface{}{
		"message": "Hello from Go!",
		"data": map[string]interface{}{
			"firstName": "Jane",
			"lastName":  "Doe",
			"age":       30,
			"email":     "jane.doe@example.com",
			"isStudent": false,
			"hobbies":   []string{"reading", "hiking", "cooking"},
			"address": map[string]string{
				"street":  "123 Main St",
				"city":    "Anytown",
				"zipCode": "12345",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func wsSender(ctx context.Context, c *Client, data interface{}) (map[string]interface{}, error) {
	c.WriteMu.Lock()
	err := c.Conn.WriteJSON(WSMessage{
		Event: "forward",
		Data:  data,
	})
	c.WriteMu.Unlock()

	if err != nil {
		return nil, err
	}

	for {
		select {
		case msg, ok := <-c.Incoming:
			if !ok {
				return nil, fmt.Errorf("client disconnected")
			}
			if msg["event"] == "response" {
				return msg, nil
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("gateway timeout")
		}
	}
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ClientID string      `json:"clientID"`
		Port     string      `json:"port"`
		Method   string      `json:"method"`
		Headers  interface{} `json:"headers"`
		Body     interface{} `json:"body"`
		Path     string      `json:"path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if req.ClientID == "" {
		http.Error(w, "clientID is required", http.StatusBadRequest)
		return
	}

	clientsMu.Lock()
	client, ok := clients[req.ClientID]
	clientsMu.Unlock()

	if !ok {
		http.Error(w, "Client not connected", http.StatusNotFound)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := wsSender(ctx, client, map[string]interface{}{
		"local_port": req.Port,
		"method":     req.Method,
		"headers":    req.Headers,
		"body":       req.Body,
		"path":       req.Path,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusGatewayTimeout)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	clientID := uuid.NewString()

	client := &Client{
		ID:       clientID,
		Conn:     conn,
		Incoming: make(chan map[string]interface{}, 8),
	}

	clientsMu.Lock()
	clients[clientID] = client
	clientsMu.Unlock()

	go client.readLoop()

	client.WriteMu.Lock()
	err = client.Conn.WriteJSON(WSMessage{
		Event:    "welcome",
		ClientID: clientID,
	})
	client.WriteMu.Unlock()

	if err != nil {
		log.Println("Write error:", err)
		return
	}

	log.Println("Client connected:", clientID)
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer ws.Close()
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}

		// switch message {
		fmt.Println(string(message))
		err = ws.WriteMessage(websocket.TextMessage, []byte("Hello, World!"))
		if err != nil {
			break
		}
	}
}

func main() {
	serveStatic()

	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			helloGet(w, r)
		} else if r.Method == http.MethodPost {
			helloPost(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/send", sendHandler)
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
