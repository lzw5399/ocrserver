package main

import (
	"fmt"
	"net/http"
	"os"

	_ "bank-ocr/docs"
	"bank-ocr/global"
	_ "bank-ocr/initialize"
	"bank-ocr/router"
)

func main() {

	//init
	//os.Setenv("TESSDATA_PREFIX","C:\\Users\\yuzu\\sdk\\msys2\\mingw64\\share\\tessdata")


	r := router.InitRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("application has started, listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		global.BANK_LOGGER.Fatal(err)
	}
}
