package config

import (
	"github.com/gorilla/sessions"
	"os"
)

var DOMAIN = ":" + os.Getenv("PORT")
var API_ROOT = "http://" + DOMAIN

var Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
