package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

func handleMagicLinkVerification(codecs []securecookie.Codec) func(c *gin.Context) {
	return func(c *gin.Context) {
		csrf := c.GetHeader(CSRF_HEADER)
		var sessionData SessionData
		securecookie.DecodeMulti(MAGIC_LINK_STORE_NAME, c.Param("magic"), &sessionData, codecs...)
		if csrf != sessionData[COOKIE_PROP_CSRF] {
			c.Header("HX-Redirect", "/login")
			c.String(http.StatusForbidden, "Login unsuccessful")
			return
		}

		// Verification Successful
		// Create Cookie
		session := sessions.Default(c)
		session.Set(COOKIE_PROP_CSRF, sessionData[COOKIE_PROP_CSRF])
		session.Set(COOKIE_PROP_EMAIL, sessionData[COOKIE_PROP_EMAIL])
		session.Save()
		// c.Header("HX-Redirect", "/secure/1")
		// c.String(http.StatusSeeOther, "Login successful")

		c.Header("Content-Type", "text/html")
		c.HTML(http.StatusOK, "redirect-no-history.go.tmpl", gin.H{
			"route": "/secure/1",
		})
	}
}
