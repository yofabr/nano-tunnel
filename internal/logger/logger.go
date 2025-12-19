package logger

import "fmt"

func WelcomLogger(clientID string) {
	fmt.Println("======================================")
	fmt.Println("Welcome! You are now connected to the remote server.")
	fmt.Printf("ClientID: %s\n", clientID)
	fmt.Println("You can use this ID to forward requests to this device.")
	fmt.Println("======================================")
}

func DisconnectLogger() {
	fmt.Println("======================================")
	fmt.Println("You have been disconnected from the remote server.")
	fmt.Println("======================================")
}
