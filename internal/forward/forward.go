package forward

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type ResponseMessage struct {
	Event    string                 `json:"event"`
	ClientID string                 `json:"clientID,omitempty"`
	Message  string                 `json:"message,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

func FetchResource(c *websocket.Conn, url, method string, headers map[string]string, body map[string]interface{}) {
	start := time.Now() // ⏱ start timer

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Fatal("Error marshalling body:", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	if headers["Content-Type"] == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error making request:", err)
	}
	defer resp.Body.Close()

	duration := time.Since(start) // ⏱ stop timer

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	event := ResponseMessage{
		Event:    "response",
		ClientID: "abc123",
		Message:  "success",
		Data: map[string]interface{}{
			"status_code": resp.StatusCode,
			"body":        string(respData),
			"headers":     resp.Header,
			"time_ms":     duration.Milliseconds(), // add as ms
			"time_string": duration.String(),       // "123ms"
		},
	}

	c.WriteJSON(event)

	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("Time taken:", duration)
	fmt.Println("Response:", string(respData))
}
