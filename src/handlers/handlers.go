package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/GAdel968/GroupieTracker/src/api"
	"github.com/GAdel968/GroupieTracker/src/models"
)

var templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))

// HomeHandler gère la page d'accueil
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Page non trouvée", http.StatusNotFound)
		return
	}

	artists, err := api.GetAllArtists()
	if err != nil {
		log.Printf("Erreur lors de la récupération des artistes: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title   string
		Artists []models.Artist
		Filters map[string]interface{}
	}{
		Title:   "Groupie Tracker - Accueil",
		Artists: artists,
		Filters: map[string]interface{}{},
	}

	err = templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Printf("Erreur lors du rendu du template: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
	}
}

func LocationsAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Locations API"))
}

// ArtistHandler gère la page de détail d'un artiste
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/artist/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID d'artiste invalide", http.StatusBadRequest)
		return
	}

	artist, err := api.GetArtistByID(id)
	if err != nil {
		log.Printf("Erreur lors de la récupération de l'artiste: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	relation, err := api.GetRelation(id)
	if err != nil {
		log.Printf("Erreur lors de la récupération des relations: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	// Get locations from your API
	if _, err := api.GetLocations(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title    string
		Artist   models.Artist
		Relation models.Relation
	}{
		Title:    "Détail de l'artiste - " + artist.Name,
		Artist:   artist,
		Relation: relation,
	}

	err = templates.ExecuteTemplate(w, "artist.html", data)
	if err != nil {
		log.Printf("Erreur lors du rendu du template: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
	}
}

// SearchHandler gère les requêtes de recherche
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Paramètre de recherche manquant", http.StatusBadRequest)
		return
	}

	// Implémentation de la recherche (voir fichier filter.go)
	results, err := SearchArtists(query)
	if err != nil {
		log.Printf("Erreur lors de la recherche: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Retourne les résultats en JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// FilterHandler gère les requêtes de filtrage
func FilterHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erreur lors de l'analyse du formulaire", http.StatusBadRequest)
		return
	}

	// Récupération des paramètres de filtre
	creationDateMin := r.FormValue("creationDateMin")
	creationDateMax := r.FormValue("creationDateMax")
	firstAlbumMin := r.FormValue("firstAlbumMin")
	firstAlbumMax := r.FormValue("firstAlbumMax")
	members := r.Form["members"]
	locations := r.Form["locations"]

	// Implémentation du filtrage (voir fichier filter.go)
	filteredArtists, err := FilterArtists(creationDateMin, creationDateMax, firstAlbumMin, firstAlbumMax, members, locations)
	if err != nil {
		log.Printf("Erreur lors du filtrage: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title   string
		Artists []models.Artist
		Filters map[string]interface{}
	}{
		Title:   "Résultats filtrés",
		Artists: filteredArtists,
		Filters: map[string]interface{}{
			"creationDateMin": creationDateMin,
			"creationDateMax": creationDateMax,
			"firstAlbumMin":   firstAlbumMin,
			"firstAlbumMax":   firstAlbumMax,
			"members":         members,
			"locations":       locations,
		},
	}

	err = templates.ExecuteTemplate(w, "filtered_results.html", data)
	if err != nil {
		log.Printf("Erreur lors du rendu du template: %v", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
	}
}
