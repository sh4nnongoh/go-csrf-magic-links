package main

import (
	"crypto/rand"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

func main() {
	authKeyMagic := make([]byte, 32)
	encryptKeyMagic := make([]byte, 32)
	rand.Read(authKeyMagic)
	rand.Read(encryptKeyMagic)
	codec := securecookie.New(authKeyMagic, encryptKeyMagic)
	codec.MaxAge(60 * 60 * 15)
	codecs := []securecookie.Codec{codec}

	authKeyCookie := make([]byte, 32)
	encryptKeyCookie := make([]byte, 32)
	rand.Read(authKeyCookie)
	rand.Read(encryptKeyCookie)
	storeCookies := cookie.NewStore(authKeyCookie, encryptKeyCookie)
	storeCookies.Options(sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 60 * 3,
		HttpOnly: true,
		Secure:   false,
	})

	router := gin.Default()
	err := router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		panic(err)
	}
	router.LoadHTMLGlob("templates/*.go.tmpl")
	router.Static("/static", "./static")
	router.Use(sessions.Sessions(COOKIE_STORE_NAME, storeCookies))
	router.POST("/magic/generate", handleMagicLinkGeneration(codecs))
	router.POST("/magic/verify/:magic", handleMagicLinkVerification(codecs))
	router.GET("/magic/verify/:magic", func(c *gin.Context) {
		c.HTML(http.StatusOK, "check-auth.go.tmpl", gin.H{
			"route": "",
		})
	})
	router.GET("/secure/:id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "check-auth.go.tmpl", gin.H{
			"route": c.Request.URL.Path,
		})
	})
	router.POST("/secure/:id", handleSecure)
	router.GET("/login", MiddlewareNoCache(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.go.tmpl", gin.H{
			"csrfToken": generateCsrf(),
		})
	})
	router.Run()
}
