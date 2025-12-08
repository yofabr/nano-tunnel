package forward

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func FetchResource(url, method string, headers map[string]string, body map[string]interface{}) {
	// marshal body to JSON
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Fatal("Error marshalling body:", err)
	}
	fmt.Println("String json body:", string(bodyBytes))
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	if headers["Content-Type"] == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// set headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error making request:", err)
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("Response:", string(respData))
}
