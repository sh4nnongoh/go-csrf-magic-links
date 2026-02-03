package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	views "github.com/sh4nnongoh/go-csrf-magic-links/templates"
)

func handleSecure(c *gin.Context) {
	csrf := c.GetHeader(csrfHeader)
	session := sessions.Default(c)
	if csrf != session.Get(cookiePropCsrf) {
		c.Header("HX-Redirect", "/login")
		c.String(http.StatusSeeOther, "Invalid CSRF")
		return
	}

	if err := views.Secure(c.Param("id"), session.Get("email").(string)).Render(c.Request.Context(), c.Writer); err != nil {
		helperAddErrorToGinContext(c, err)
		return
	}
}
