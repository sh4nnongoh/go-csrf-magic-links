package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func handleSecure(c *gin.Context) {
	csrf := c.GetHeader(csrfHeader)
	session := sessions.Default(c)
	if csrf != session.Get(cookiePropCsrf) {
		c.Header("HX-Redirect", "/login")
		c.String(http.StatusSeeOther, "Invalid CSRF")
		return
	}
	c.HTML(http.StatusOK, "secure.go.tmpl", gin.H{
		"count": c.Param("id"),
		"email": session.Get("email"),
	})
}
