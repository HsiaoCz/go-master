package main

import (
	"net/http"

	"github.com/uptrace/bunrouter"
)

func main() {
	router := bunrouter.New()
	router.GET("/", func(w http.ResponseWriter, req bunrouter.Request) error {
		return bunrouter.JSON(w, map[string]any{
			"message": "Hello, World!",
		})
	})
	http.ListenAndServe(":8080", router)
}
