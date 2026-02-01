package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gorilla/securecookie"
	gsessions "github.com/gorilla/sessions"
)

func NewJsonStore(codecs []securecookie.Codec) sessions.Store {
	cs := &gsessions.CookieStore{
		Codecs: codecs,
		Options: &gsessions.Options{
			Path:     "/",
			MaxAge:   86400 * 30,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		},
	}
	cs.MaxAge(cs.Options.MaxAge)

	return &store{cs}
}

type store struct {
	*gsessions.CookieStore
}

func (c *store) Options(options sessions.Options) {
	c.CookieStore.Options = options.ToGorillaOptions()
}
