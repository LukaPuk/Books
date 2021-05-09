package main

import (
	"fmt"
	"github.com/LukaPuk/Books/internal/driver"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
	"os"
)

func init() {
	gotenv.Load() // loads env file

}

// add neccessary json format or error

func main() {

	err := driver.InitPostgres()
	defer driver.DB.Close()
	if err != nil {
		log.Fatal()
	}

	fmt.Println("test")

	//_,  err = render.CreateTemplateCache()
	//if err != nil {
	//	log.Fatal("cannot create template cache")
	//}

	router := InitRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}

}
