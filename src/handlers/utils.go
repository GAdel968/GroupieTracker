package handlers

import (
	"html/template"
)

// Fonction utilitaire pour vérifier si un tableau contient une valeur
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// FuncMap définit les fonctions personnalisées pour les templates
var funcMap = template.FuncMap{
	"contains": containsString,
}

// Mise à jour de la variable templates avec les fonctions personnalisées
func init() {
	templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))
}
