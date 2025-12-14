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
	Event   string      `json:"event"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

var (
	clients   = make(map[string]*websocket.Conn)
	clientsMu sync.Mutex
)

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

func wsSender(ctx context.Context, conn *websocket.Conn, data interface{}) (map[string]interface{}, error) {
	request := WSMessage{
		Event: "forward",
		Data:  data,
	}

	if err := conn.WriteJSON(request); err != nil {
		return nil, err
	}

	responseChan := make(chan map[string]interface{}, 1)

	go func() {
		for {
			var msg map[string]interface{}
			if err := conn.ReadJSON(&msg); err != nil {
				return
			}
			if msg["event"] == "response" {
				responseChan <- msg
				return
			}
		}
	}()

	select {
	case resp := <-responseChan:
		return resp, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("gateway timeout")
	}
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ClientID string      `json:"clientID"`
		Port     int         `json:"port"`
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

	clientsMu.Lock()
	clients[clientID] = conn
	clientsMu.Unlock()

	// Welcome message
	conn.WriteJSON(WSMessage{
		Event:   "welcome",
		Message: "",
		Data: map[string]string{
			"clientID": clientID,
		},
	})

	log.Println("Client connected:", clientID)

	defer func() {
		clientsMu.Lock()
		delete(clients, clientID)
		clientsMu.Unlock()
		conn.Close()
		log.Println("Client disconnected:", clientID)
	}()
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
