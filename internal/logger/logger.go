package logger

import "fmt"

func WelcomLogger(clientID string) {
	fmt.Println("======================================")
	fmt.Println("Welcome! You are now connected to the server.")
	fmt.Printf("ClientID: %s\n", clientID)
	fmt.Println("You can use this ID to forward requests to this device.")
	fmt.Println("======================================")
}
