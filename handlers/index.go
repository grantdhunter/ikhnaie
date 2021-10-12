package handlers

import (
	"log"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	log.Print(r.URL)
	accessToken := r.URL.Query()["access_token"]
	log.Print(accessToken)
	if len(accessToken) != 0 {
		http.Redirect(w, r, "/track/", http.StatusFound)
	} else {
		http.Redirect(w, r, "/auth/", http.StatusFound)
	}
}
