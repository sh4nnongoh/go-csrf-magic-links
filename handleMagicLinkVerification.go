package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

func handleMagicLinkVerification(codecs []securecookie.Codec) func(c *gin.Context) {
	return func(c *gin.Context) {
		csrf := c.GetHeader(csrfHeader)
		var sessionData SessionData

		if err := securecookie.DecodeMulti(magicLinkStoreName, c.Param("magic"), &sessionData, codecs...); err != nil {
			_ = fmt.Errorf("failed to decode session: %w", err)
		}
		if csrf != sessionData[cookiePropCsrf] {
			c.Header("HX-Redirect", "/login")
			c.String(http.StatusForbidden, "Login unsuccessful")
			return
		}

		// Verification Successful
		// Create Cookie
		session := sessions.Default(c)
		session.Set(cookiePropCsrf, sessionData[cookiePropCsrf])
		session.Set(cookiePropEmail, sessionData[cookiePropEmail])
		if err := session.Save(); err != nil {
			_ = fmt.Errorf("failed to save session: %w", err)
		}
		c.Header("Content-Type", "text/html")
		c.HTML(http.StatusOK, "redirect-no-history.go.tmpl", gin.H{
			"route": "/secure/1",
		})
	}
}
