package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type token struct {
	Token string `json:"token"`
	Misc  string `json:"misc"`
}

var testToken = token{Token: "123456abc", Misc: "hello, this a token"}
var testLog string = "hello, this is a test log"

func main() {
	router := gin.Default()
	router.GET("/token", inviteTokenGen)
	router.Run("localhost:8080")
}

func inviteTokenGen(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, testToken)
	fmt.Printf(testLog)
}
