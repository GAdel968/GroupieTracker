package main

import (
	"log"
	"net/http"

	"Groupie-Tracker/src/handlers"
)

func main() {
	// Configuration des fichiers statiques
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Configuration des routes
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/artist/", handlers.ArtistHandler)
	http.HandleFunc("/search", handlers.SearchHandler)
	http.HandleFunc("/filter", handlers.FilterHandler)

	// Ajout d'une route API pour les locations
	http.HandleFunc("/api/locations", handlers.LocationsAPIHandler)

	// Démarrage du serveur
	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
