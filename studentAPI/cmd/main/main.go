package main

import (
	"studentAPI/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	handler.SetupHandler(r)
	r.Run(":8080")
}
