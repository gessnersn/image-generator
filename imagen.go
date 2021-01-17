package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var baseTpl = template.Must(template.ParseFiles("base.html"))

func main() {
	port := getPortFromEnv()
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/imagen/", imagenHandler)
	log.Fatal(http.ListenAndServe(":" + port, mux))
}

func getPortFromEnv() string {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return port
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(
		writer,
		"This is the index page, you can also explore the '/imagen' subpage!",
	)
}

func imagenHandler(writer http.ResponseWriter, request *http.Request) {
	data := struct {
		Title string
		Image string
	}{
		Title: "Imagen - Image Generator",
		Image: "/assets/black640x480.jpg",
	}
	err := baseTpl.Execute(writer, data)
	check(err)
}

func check (err error) {
	if err != nil {
		log.Fatal(err)
	}
}
