package main

import (
	"bank-ocr/router"
	"log"
	"net/http"
	"os"
)



func main() {
	r := router.InitRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
		//logger.Fatalln("Required env `PORT` is not specified.")
	}

	logger.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Println(err)
	}
}
