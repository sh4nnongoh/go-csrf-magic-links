package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"

	//nolint:gosec // G108: pprof is only exposed on internal environments
	_ "net/http/pprof"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

	authKeyMagic := make([]byte, 32)
	encryptKeyMagic := make([]byte, 32)
	if _, err := rand.Read(authKeyMagic); err != nil {
		_ = fmt.Errorf("failed to generate authKeyMagic: %w", err)
	}
	if _, err := rand.Read(encryptKeyMagic); err != nil {
		_ = fmt.Errorf("failed to generate encryptKeyMagic: %w", err)
	}
	codec := securecookie.New(authKeyMagic, encryptKeyMagic)
	codec.MaxAge(60 * 60 * 15)
	codecs := []securecookie.Codec{codec}

	authKeyCookie := make([]byte, 32)
	encryptKeyCookie := make([]byte, 32)
	if _, err := rand.Read(authKeyMagic); err != nil {
		_ = fmt.Errorf("failed to generate authKeyCookie: %w", err)
	}
	if _, err := rand.Read(encryptKeyMagic); err != nil {
		_ = fmt.Errorf("failed to generate encryptKeyCookie: %w", err)
	}
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
