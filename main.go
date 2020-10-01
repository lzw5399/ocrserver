package main

import (
	"bank-ocr/global"
	"fmt"
	"net/http"
	"os"

	_ "bank-ocr/docs"
	_ "bank-ocr/initialize"
	"bank-ocr/router"
)

func main() {
	r := router.InitRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		global.BANK_LOGGER.Fatal(err)
	}
}
