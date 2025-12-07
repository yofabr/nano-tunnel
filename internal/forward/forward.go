package forward

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func FetchResource(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Couldn't find resource", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	fmt.Println("Data:", string(data))
}
