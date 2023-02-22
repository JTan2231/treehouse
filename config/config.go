package config

import (
	"github.com/gorilla/sessions"
	"os"
)

const DOMAIN = "0.0.0.0:8080"
const API_ROOT = "http://" + DOMAIN

var Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
