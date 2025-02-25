package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/GAdel968/GroupieTracker/src/api"
	"github.com/GAdel968/GroupieTracker/src/models"
)

// SearchResult représente un résultat de recherche avec son type
type SearchResult struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

// SearchArtists recherche les artistes correspondant à la requête
func SearchArtists(query string) ([]SearchResult, error) {
	var results []SearchResult
	query = strings.ToLower(query)

	artists, err := api.GetAllArtists()
	if err != nil {
		return nil, err
	}

	for _, artist := range artists {
		// Recherche par nom d'artiste
		if strings.Contains(strings.ToLower(artist.Name), query) {
			results = append(results, SearchResult{
				Value: artist.Name,
				Type:  "artiste/groupe",
			})
		}

		// Recherche par membres
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
				results = append(results, SearchResult{
					Value: member,
					Type:  "membre",
				})
			}
		}

		// Recherche par date de création
		creationDateStr := strconv.Itoa(artist.CreationDate)
		if strings.Contains(creationDateStr, query) {
			results = append(results, SearchResult{
				Value: creationDateStr,
				Type:  "date de création",
			})
		}

		// Recherche par date du premier album
		if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
			results = append(results, SearchResult{
				Value: artist.FirstAlbum,
				Type:  "date du premier album",
			})
		}
	}

	// Recherche par lieux de concert
	locations, err := api.GetAllLocations()
	if err != nil {
		return nil, err
	}

	for _, locationList := range locations.Locations {
		for _, location := range locationList {
			if strings.Contains(strings.ToLower(location), query) {
				results = append(results, SearchResult{
					Value: location,
					Type:  "lieu de concert",
				})
			}
		}
	}

	return results, nil
}

// FilterArtists filtre les artistes selon les critères spécifiés
func FilterArtists(creationDateMin, creationDateMax, firstAlbumMin, firstAlbumMax string, membersFilter, locationsFilter []string) ([]models.Artist, error) {
	artists, err := api.GetAllArtists()
	if err != nil {
		return nil, err
	}

	var filteredArtists []models.Artist

	for _, artist := range artists {
		// Vérifie si l'artiste correspond à tous les filtres
		if matchesCreationDateFilter(artist, creationDateMin, creationDateMax) &&
			matchesFirstAlbumFilter(artist, firstAlbumMin, firstAlbumMax) &&
			matchesMembersFilter(artist, membersFilter) &&
			matchesLocationsFilter(artist, locationsFilter) {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	return filteredArtists, nil
}

// Vérifie si l'artiste correspond au filtre de date de création
func matchesCreationDateFilter(artist models.Artist, minStr, maxStr string) bool {
	// Si aucun filtre n'est spécifié, l'artiste correspond
	if minStr == "" && maxStr == "" {
		return true
	}

	min, err := strconv.Atoi(minStr)
	if err != nil || minStr == "" {
		min = 0
	}

	max, err := strconv.Atoi(maxStr)
	if err != nil || maxStr == "" {
		max = 9999
	}

	return artist.CreationDate >= min && artist.CreationDate <= max
}

// Vérifie si l'artiste correspond au filtre de date du premier album
func matchesFirstAlbumFilter(artist models.Artist, minStr, maxStr string) bool {
	// Si aucun filtre n'est spécifié, l'artiste correspond
	if minStr == "" && maxStr == "" {
		return true
	}

	albumDate, err := time.Parse("02-01-2006", artist.FirstAlbum)
	if err != nil {
		return false
	}

	var minDate, maxDate time.Time

	if minStr != "" {
		minDate, err = time.Parse("2006-01-02", minStr)
		if err != nil {
			minDate = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
		}
	} else {
		minDate = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	if maxStr != "" {
		maxDate, err = time.Parse("2006-01-02", maxStr)
		if err != nil {
			maxDate = time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)
		}
	} else {
		maxDate = time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)
	}

	return albumDate.After(minDate) && albumDate.Before(maxDate)
}

// Vérifie si l'artiste correspond au filtre de nombre de membres
func matchesMembersFilter(artist models.Artist, membersFilter []string) bool {
	// Si aucun filtre n'est spécifié, l'artiste correspond
	if len(membersFilter) == 0 || (len(membersFilter) == 1 && membersFilter[0] == "") {
		return true
	}

	// Convertit les filtres en nombres
	var memberCounts []int
	for _, countStr := range membersFilter {
		count, err := strconv.Atoi(countStr)
		if err == nil {
			memberCounts = append(memberCounts, count)
		}
	}

	// Vérifie si le nombre de membres de l'artiste correspond à l'un des filtres
	for _, count := range memberCounts {
		if len(artist.Members) == count {
			return true
		}
	}

	return false
}

// Vérifie si l'artiste correspond au filtre de lieux
func matchesLocationsFilter(artist models.Artist, locationsFilter []string) bool {
	// Si aucun filtre n'est spécifié, l'artiste correspond
	if len(locationsFilter) == 0 || (len(locationsFilter) == 1 && locationsFilter[0] == "") {
		return true
	}

	// Récupération des relations pour obtenir les lieux de concert
	relation, err := api.GetRelation(artist.ID)
	if err != nil {
		return false
	}

	// Vérifie si l'un des lieux de l'artiste correspond à l'un des filtres
	for location := range relation.DatesLocations {
		locationLower := strings.ToLower(location)
		for _, filter := range locationsFilter {
			if strings.Contains(locationLower, strings.ToLower(filter)) {
				return true
			}
		}
	}

	return false
}
