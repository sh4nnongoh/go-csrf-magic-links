package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func helperAddErrorToGinContext(c *gin.Context, err error) {
	if ginError := c.Error(err); ginError != nil {
		log.Println("failed to attach error to gin context: ", ginError)
	}
}
