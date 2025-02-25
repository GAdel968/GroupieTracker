package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/GAdel968/GroupieTracker/src/models"
)

const (
	ArtistsURL   = "https://groupietrackers.herokuapp.com/api/artists"
	LocationsURL = "https://groupietrackers.herokuapp.com/api/locations"
	DatesURL     = "https://groupietrackers.herokuapp.com/api/dates"
	RelationURL  = "https://groupietrackers.herokuapp.com/api/relation"
)

// GetAllArtists récupère tous les artistes depuis l'API
func GetAllArtists() ([]models.Artist, error) {
	resp, err := http.Get(ArtistsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var artists []models.Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

// GetArtistByID récupère un artiste spécifique par son ID
func GetArtistByID(id int) (models.Artist, error) {
	url := fmt.Sprintf("%s/%d", ArtistsURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return models.Artist{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.Artist{}, err
	}

	var artist models.Artist
	err = json.Unmarshal(body, &artist)
	if err != nil {
		return models.Artist{}, err
	}

	return artist, nil
}

// GetRelation récupère les relations pour un artiste par son ID
func GetRelation(id int) (models.Relation, error) {
	url := fmt.Sprintf("%s/%d", RelationURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return models.Relation{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.Relation{}, err
	}

	var relation models.Relation
	err = json.Unmarshal(body, &relation)
	if err != nil {
		return models.Relation{}, err
	}

	return relation, nil
}

// GetAllLocations récupère toutes les informations de localisation
func GetAllLocations() (models.Location, error) {
	resp, err := http.Get(LocationsURL)
	if err != nil {
		return models.Location{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.Location{}, err
	}

	var location models.Location
	err = json.Unmarshal(body, &location)
	if err != nil {
		return models.Location{}, err
	}

	return location, nil
}

// GetAllDates récupère toutes les informations de dates
func GetAllDates() (models.Date, error) {
	resp, err := http.Get(DatesURL)
	if err != nil {
		return models.Date{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.Date{}, err
	}

	var date models.Date
	err = json.Unmarshal(body, &date)
	if err != nil {
		return models.Date{}, err
	}

	return date, nil
}

// GetLocations retrieves the locations from the API
func GetLocations() ([]string, error) {
	// Implement the logic to get locations
	// For now, return a dummy slice of locations and nil error
	return []string{"Location1", "Location2"}, nil
}
