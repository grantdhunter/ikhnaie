package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/grantdhunter/ikhnaie/models"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/auth.html")
	t.Execute(w, models.Config)
}

var httpClient = http.Client{}

func OauthRedirectHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Printf("could not parse query: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	code := r.FormValue("code")

	reqURL := fmt.Sprintf("%s/access_token?client_id=%s&client_secret=%s&code=%s",
		models.Config.Github_oauth_url,
		models.Config.Github_oauth_clientid,
		models.Config.Github_oauth_client_secret,
		code)

	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		log.Printf("could not create HTTP request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	// We set this header since we want the response
	// as JSON
	req.Header.Set("Accept", "application/json")

	// Send out the HTTP request
	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("could not send HTTP request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer res.Body.Close()

	// Parse the request body into the `OAuthAccessResponse` struct
	var t models.OAuthAccessResponse
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		log.Printf("could not parse JSON response: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	// Finally, send a response to redirect the user to the "welcome" page
	// with the access token
	w.Header().Set("Location", fmt.Sprintf("/app/?token=%s", t.AccessToken))
	w.Header().Set("Authentication", fmt.Sprintf("Bearer %s", t.AccessToken))
	w.Header().Set("WWW-Authenticate", "Bearer realm=\"All\"")
	w.WriteHeader(http.StatusFound)
}
