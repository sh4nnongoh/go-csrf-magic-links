package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	views "github.com/sh4nnongoh/go-csrf-magic-links/templates"
)

func handleMagicLinkGeneration(codecs []securecookie.Codec) func(c *gin.Context) {
	return func(c *gin.Context) {
		csrf := c.GetHeader("X-CSRF-Token")
		encoded, _ := securecookie.EncodeMulti(
			magicLinkStoreName,
			SessionData{
				"csrf":  csrf,
				"email": c.PostForm("email"),
			},
			codecs...,
		)
		views.MagicGenerate("http://localhost:8080/magic/verify/"+encoded).Render(c.Request.Context(), c.Writer)
	}
}
