package main

import (
	"log"
	"loupgarou/internal/handler"
	"loupgarou/internal/service"
	"loupgarou/internal/storage"
	"net/http"
)

func main() {
	repo := storage.NewFileGameRepository("games")
	svc := service.NewGameService(repo)
	h := handler.NewGameHandler(svc)

	http.HandleFunc("/game", h.HandleGame)
	http.HandleFunc("/status", h.HandleStatus)
	http.HandleFunc("/step", h.HandleStep)
	http.Handle("/", http.FileServer(http.Dir("web")))

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
