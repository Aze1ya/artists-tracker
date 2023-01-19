package main

import (
	"fmt"
	"log"
	"net/http"

	handlers "alem/cmd"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/artist/", handlers.Artistdata)
	mux.HandleFunc("/filters/", handlers.Filterdata)

	fileServer := http.FileServer(http.Dir("./ui/style/"))
	mux.Handle("/style/", http.StripPrefix("/style", fileServer))

	handlers.NewClient()

	fmt.Println("http://localhost:8080/")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
