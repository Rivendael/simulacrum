package main

import (
	"fmt"
	"log"

	"simulacrum/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/obscure", handlers.HandleObscure)

	fmt.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
