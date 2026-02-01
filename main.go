package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	//nolint:gosec // G108: pprof is only exposed on internal environments
	_ "net/http/pprof"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

//nolint:funlen
func main() {
	// Start pprof HTTP server
	go func() {
		srv := &http.Server{
			Addr:         ":6060",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  5 * time.Second,
		}

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	// For Magic Links
	authKeyMagic := securecookie.GenerateRandomKey(32)
	encryptKeyMagic := securecookie.GenerateRandomKey(32)
	codec := securecookie.New(authKeyMagic, encryptKeyMagic)
	codec.SetSerializer(JSONEncoder{})
	codec.MaxAge(60 * 60 * 15)
	codecs := []securecookie.Codec{codec}

	// For Cookie Store
	authKeyCookie := securecookie.GenerateRandomKey(32)
	encryptKeyCookie := securecookie.GenerateRandomKey(32)
	codecCookie := securecookie.New(authKeyCookie, encryptKeyCookie)
	codecCookie.SetSerializer(JSONEncoder{})
	codecCookie.MaxAge(60 * 60 * 60 * 3)
	codecsCookie := []securecookie.Codec{codecCookie}
	storeCookies := NewJsonStore(codecsCookie)
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
	router.Use(sessions.Sessions(cookieStoreName, storeCookies))
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

	if err := router.Run(); err != nil {
		_ = fmt.Errorf("failed to run router: %w", err)
	}
}
