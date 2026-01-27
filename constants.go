// This file contains the common constants used throughout the webserver lifecycle.
package main

const (
	csrfHeader         = "x-csrf-token"
	magicLinkStoreName = "magic"
	cookieStoreName    = "cookies"
	cookiePropCsrf     = "csrf"
	cookiePropEmail    = "email"
)
