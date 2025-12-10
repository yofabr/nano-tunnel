package forward

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type ResponseMessage struct {
	Event       string                 `json:"event"`
	ClientID    string                 `json:"clientID,omitempty"`
	Message     string                 `json:"message,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Status_Code int                    `json:"status_code,omitempty"`
	Headers     http.Header            `json:"headers,omitempty"`
	TimeString  string                 `json:"time_string,omitempty"`
	TimeMs      int64                  `json:"time_ms,omitempty"`
}

func FetchResource(c *websocket.Conn, clientID, url, method string, headers map[string]string, body map[string]interface{}) {
	start := time.Now()

	if method == "" {
		method = http.MethodGet
	}
	method = strings.ToUpper(method)

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		sendError(c, clientID, fmt.Errorf("marshal body: %w", err), start)
		return
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		sendError(c, clientID, fmt.Errorf("create request: %w", err), start)
		return
	}

	if headers == nil {
		headers = map[string]string{}
	}

	if headers["Content-Type"] == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		sendError(c, clientID, fmt.Errorf("forward request: %w", err), start)
		return
	}
	defer resp.Body.Close()

	duration := time.Since(start)

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		sendError(c, clientID, fmt.Errorf("read response: %w", err), start)
		return
	}

	event := ResponseMessage{
		Event:    "response",
		ClientID: clientID,
		Message:  "success",
		Data: map[string]interface{}{
			"response": string(respData),
		},
		Status_Code: resp.StatusCode,
		Headers:     resp.Header,
		TimeString:  duration.String(),
		TimeMs:      duration.Milliseconds(),
	}

	if err := c.WriteJSON(event); err != nil {
		log.Printf("failed to write response to websocket: %v", err)
	}
}

func sendError(c *websocket.Conn, clientID string, err error, start time.Time) {
	duration := time.Since(start)

	event := ResponseMessage{
		Event:    "response",
		ClientID: clientID,
		Message:  "error",
		Data: map[string]interface{}{
			"error": err.Error(),
		},
		Status_Code: http.StatusBadGateway,
		TimeString:  duration.String(),
		TimeMs:      duration.Milliseconds(),
	}

	if writeErr := c.WriteJSON(event); writeErr != nil {
		log.Printf("failed to write error response to websocket: %v (original: %v)", writeErr, err)
	}
}
