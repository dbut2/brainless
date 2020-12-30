package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/puzzle-solvers/brainless/pkg/binairo"
	"github.com/puzzle-solvers/brainless/pkg/pipes"
	"github.com/puzzle-solvers/brainless/pkg/sudoku"
	"net/http"
)

func main() {
	router := chi.NewRouter()

	router.HandleFunc("/binairo", binairo.Solve)
	router.HandleFunc("/pipes", pipes.Solve)
	router.HandleFunc("/sudoku", sudoku.Solve)

	fmt.Println("Starting server...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Shutting down server")
}
