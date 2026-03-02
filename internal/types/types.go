package types

type WsData struct {
	LocalPort string                 `json:"local_port,omitempty"`
	Path      string                 `json:"path,omitempty"`
	Method    string                 `json:"method,omitempty"`
	Headers   map[string]string      `json:"headers,omitempty"`
	Body      map[string]interface{} `json:"body,omitempty"`
}

type Message struct {
	Event    string `json:"event"`
	ClientID string `json:"clientID,omitempty"`
	Message  string `json:"message,omitempty"`
	Data     WsData `json:"data,omitempty"`
}

type ResponseMessage struct {
	Event      string                 `json:"event"`
	ClientID   string                 `json:"clientID,omitempty"`
	Message    string                 `json:"message,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
	StatusCode int                    `json:"status_code,omitempty"`
	Headers    map[string][]string    `json:"headers,omitempty"`
	TimeString string                 `json:"time_string,omitempty"`
	TimeMs     int64                  `json:"time_ms,omitempty"`
}
