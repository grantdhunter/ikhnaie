package handlers

import (
	"html/template"
	"net/http"

	"github.com/grantdhunter/ikhnaie/models"
)

func TrackHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/track.html")
	t.Execute(w, models.Config)
}
