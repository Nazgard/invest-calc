package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Printf("Server starting at http://localhost%s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-stop
	fmt.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	fmt.Println("Server stopped")
}
