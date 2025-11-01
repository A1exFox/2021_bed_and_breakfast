package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a1exfox/go-course/pkg/config"
	"github.com/a1exfox/go-course/pkg/handlers"
	"github.com/a1exfox/go-course/pkg/render"
)

const portNumber = ":8080"

func main() {

	var app config.AppConfig
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache ", err)
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	http.HandleFunc("/", handlers.Repo.Home)
	http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println("Starting application on port", portNumber)
	http.ListenAndServe(portNumber, nil)
}
