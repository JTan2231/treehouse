package config

import (
    //"net/http"
	//"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"os"
)
const DOMAIN = "localhost:8080"
const API_ROOT = "http://" + DOMAIN

var Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))