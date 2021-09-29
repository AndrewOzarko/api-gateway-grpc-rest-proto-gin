package main

import (
	"context"
	"gateway/gateway/client"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	timeout     = time.Second
	user_client = client.UserClient{}
)

func UserRegister(r *gin.RouterGroup) {
	r.GET("/", GetUsers)
}

func GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	data, err := user_client.GetUsers(&ctx)
	response(c, data, err)
}

func response(c *gin.Context, data interface{}, err error) {
	statusCode := http.StatusOK
	var errorMessage string
	if err != nil {
		log.Println("Server Error Occured:", err)
		errorMessage = strings.Title(err.Error())
		statusCode = http.StatusInternalServerError
	}
	c.JSON(statusCode, gin.H{"data": data, "error": errorMessage})
}

func main() {
	log.Println("Gateway Service")

	r := gin.Default()
	r.Use(cors.Default())

	api := r.Group("/api")
	v1 := api.Group("/v1")

	UserRegister(v1.Group("/users"))
	r.Run()
}
