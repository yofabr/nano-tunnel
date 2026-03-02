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
	"github.com/yofabr/nano-tunnel/internal/types"
)

var validMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}

func FetchResource(c *websocket.Conn, clientID, targetURL, method string, headers map[string]string, body map[string]interface{}) {
	start := time.Now()

	if method == "" {
		method = http.MethodGet
	}
	method = strings.ToUpper(method)

	isValidMethod := false
	for _, m := range validMethods {
		if method == m {
			isValidMethod = true
			break
		}
	}

	if !isValidMethod {
		sendError(c, clientID, fmt.Errorf("invalid HTTP method: %s", method), start)
		return
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		sendError(c, clientID, fmt.Errorf("marshal body: %w", err), start)
		return
	}

	req, err := http.NewRequest(method, targetURL, bytes.NewBuffer(bodyBytes))
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

	event := types.ResponseMessage{
		Event:      "response",
		ClientID:   clientID,
		Message:    "success",
		Data:       map[string]interface{}{"response": string(respData)},
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		TimeString: duration.String(),
		TimeMs:     duration.Milliseconds(),
	}

	if err := c.WriteJSON(event); err != nil {
		log.Printf("failed to write response to websocket: %v", err)
	}
}

func sendError(c *websocket.Conn, clientID string, err error, start time.Time) {
	duration := time.Since(start)

	event := types.ResponseMessage{
		Event:      "response",
		ClientID:   clientID,
		Message:    "error",
		Data:       map[string]interface{}{"error": err.Error()},
		StatusCode: http.StatusBadGateway,
		TimeString: duration.String(),
		TimeMs:     duration.Milliseconds(),
	}

	if writeErr := c.WriteJSON(event); writeErr != nil {
		log.Printf("failed to write error response to websocket: %v (original: %v)", writeErr, err)
	}
}
