// cmd/main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"MapleTools/internal/handlers" // adjust import path
)

func main() {
	http.HandleFunc("/", handlers.FormHandler)

	fmt.Println("Server starting at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
