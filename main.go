/*
@Title           Blog Platform API
@Description     This is the main entry point for the Blog Platform API server.
@Version         1.0
@Host            localhost:5000
@BasePath        /
@Schemes         http
*/
package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/Aman913k/docs"

	"github.com/Aman913k/routes"
)

func main() {
	r := routes.Router()
	fmt.Println("Server is getting Started...")

	log.Fatal(http.ListenAndServe(":5000", r))
	fmt.Println("Listening at port 5000...")
}
