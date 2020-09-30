package main

import (
	"net/http"
	"os"

	"bank-ocr/router"
)

func main() {
	r := router.InitRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	//global.BANK_LOGGER.Info("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		//global.BANK_LOGGER.Fatal(err)
	}
}
