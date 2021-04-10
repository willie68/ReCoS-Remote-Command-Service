package handler

import (
	"encoding/base64"
	"net/http"
	"strings"
)

var AuthenticationConfig = AuthConfig{
	Username: "admin",
	Password: "recosadmin",
}

type AuthConfig struct {
	Username string
	Password string
}

// AuthCheck implements a simple middleware handler for adding basic http auth to a route.
func AuthCheck() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if AuthenticationConfig.Password != "" {

				auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

				if len(auth) != 2 || auth[0] != "Basic" {
					http.Error(w, "{\"message\": \"authorization failed\"}", http.StatusUnauthorized)
					return
				}

				payload, _ := base64.StdEncoding.DecodeString(auth[1])
				pair := strings.SplitN(string(payload), ":", 2)

				if len(pair) != 2 || !validate(pair[0], pair[1]) {
					http.Error(w, "{\"message\": \"authorization failed\"}", http.StatusUnauthorized)
					return
				}
			}

			next.ServeHTTP(w, r)
			return
		})
	}
}

func validate(username, password string) bool {
	if username == AuthenticationConfig.Username && password == AuthenticationConfig.Password {
		return true
	}
	return false
}
