package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

func handleMagicLinkGeneration(codecs []securecookie.Codec) func(c *gin.Context) {
	return func(c *gin.Context) {
		csrf := c.GetHeader("X-CSRF-Token")
		encoded, _ := securecookie.EncodeMulti(
			MAGIC_LINK_STORE_NAME,
			SessionData{
				"csrf":  csrf,
				"email": c.PostForm("email"),
			},
			codecs...,
		)
		c.HTML(http.StatusOK, "magic-generate.go.tmpl", gin.H{
			"magic": "http://localhost:8080/magic/verify/" + encoded,
		})
	}
}
