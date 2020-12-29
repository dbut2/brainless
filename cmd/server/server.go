package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/puzzle-solvers/brainless/internal/pipes"
	"net/http"
)

func main() {
	router := chi.NewRouter()
	router.HandleFunc("/solve", pipes.Solve)
	fmt.Println("Starting server...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Printf("Sever failed to start: %s\n", err.Error())
		return
	}
	fmt.Println("Server stopped.")
}
