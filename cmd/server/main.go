package main

import (
	"fmt"
	"log"
	"net/http"

	"invest-calc/internal/infrastructure"
	httpHandler "invest-calc/internal/interfaces/http"
	"invest-calc/internal/usecases"
)

func main() {
	store := infrastructure.NewMemoryStore()
	calc := usecases.NewCalculator(store)

	handler, err := httpHandler.NewHandler(calc)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	addr := ":8080"
	fmt.Printf("Server starting at http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
