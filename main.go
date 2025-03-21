package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var (
	templates      = template.Must(template.ParseGlob("templates/*.html"))
	filterTemplate = template.Must(template.ParseFiles("templates/filter.html"))
	artistTemplate = template.Must(template.ParseFiles("templates/artists.html"))
)

type Artists struct {
	ID           int      `json:"id"`
	IMAGE        string   `json:"image"`
	NAME         string   `json:"name"`
	MEMBERS      []string `json:"members"`
	CREA_DATE    int      `json:"creationDate"`
	FIRST_ALBUM  string   `json:"firstAlbum"`
	LOCATIONS    string   `json:"locations"`
	CONCERT_DATE string   `json:"concertDates"`
	RELATIONS    string   `json:"relations"`
}

type Datas struct {
	ArtistsData   []Artists   `json:"artists"`
	LocationsData []Locations `json:"locations"`
	RelationData  []Relations `json:"relations"`
	DatesData     []Dates     `json:"dates"`
}

type Locations struct {
	ID        int      `json:"id"`
	LOCATIONS []string `json:"locations"`
	DATES     string   `json:"dates"`
}

type Dates struct {
	ID    int      `json:"id"`
	DATES []string `json:"dates"`
}

type Relations struct {
	ID         int                 `json:"id"`
	DATESLOCAT map[string][]string `json:"datesLocations"`
}

func main() {
	// Configuration des fichiers statiques
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/artists/", artistHandler)
	http.HandleFunc("/page", pageHandler)
	http.HandleFunc("/filter", filterPageHandler)

	// Démarrage du serveur
	fmt.Println("Serveur démarré sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	data := loadData()
	if err := templates.ExecuteTemplate(w, "page.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	data := loadData()
	// Trouver l'artiste correspondant à l'ID
	// ...

	if err := templates.ExecuteTemplate(w, "artists.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func filterPageHandler(w http.ResponseWriter, r *http.Request) {
	data := loadData()
	// Logique de filtrage ici
	// ...

	if err := templates.ExecuteTemplate(w, "filter.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loadData() Datas {
	var data Datas
	files := []string{"artists.json", "locations.json", "dates.json", "relations.json"}

	for _, file := range files {
		jsonFile, err := os.Open(filepath.Join("data", file))
		if err != nil {
			log.Printf("Erreur lors de l'ouverture de %s: %v", file, err)
			continue
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		switch file {
		case "artists.json":
			json.Unmarshal(byteValue, &data.ArtistsData)
		case "locations.json":
			json.Unmarshal(byteValue, &data.LocationsData)
		case "dates.json":
			json.Unmarshal(byteValue, &data.DatesData)
		case "relations.json":
			json.Unmarshal(byteValue, &data.RelationData)
		}
	}

	return data
}

/*--------------------------------------------------------------------------------------------
------------------------------- Fonction Gestionnaire ArtistPage -----------------------------
----------------------------------------------------------------------------------------------*/

func detailedArtistHandler(w http.ResponseWriter, r *http.Request) {
	// ART variable de type structure DATAS
	ART := loadData()

	// Récupère l'ID de l'Artiste
	id := r.URL.Path[9:]
	p, err := strconv.Atoi(id)
	if err != nil || p < 1 || p > len(ART.ArtistsData) {
		http.Error(w, "Artiste non trouvé", http.StatusNotFound)
		return
	}

	var infoArt = make(map[string]any)

	// Stocke les informations utiles dans la carte infoArt
	infoArt["IMAGE"] = ART.ArtistsData[p-1].IMAGE
	infoArt["NAME"] = ART.ArtistsData[p-1].NAME
	infoArt["MEMBERS"] = ART.ArtistsData[p-1].MEMBERS
	infoArt["CREA_DATE"] = ART.ArtistsData[p-1].CREA_DATE
	infoArt["FIRST_ALBUM"] = ART.ArtistsData[p-1].FIRST_ALBUM
	infoArt["DATESLOCAT"] = ART.RelationData[p-1].DATESLOCAT

	if error := artistTemplate.Execute(w, infoArt); error != nil {
		log.Fatal(error)
	}
}

/*--------------------------------------------------------------------------------------------
------------------------------- Fonctions pour la partie Filtre ------------------------------
----------------------------------------------------------------------------------------------*/

// la fonction getFilters récupère, à partir des filtres sur la page HTML Page, les différentes données et les stocke dans une carte (filters)
func getFilters(r *http.Request) map[string]string {
	filters := make(map[string]string)

	filters["members"] = r.URL.Query().Get("members")
	filters["firstAlbum"] = r.URL.Query().Get("firstAlbum")
	filters["creationDate"] = r.URL.Query().Get("creationDate")
	filters["CitySearch"] = r.URL.Query().Get("CitySearch")

	return filters
}

// /filter.html
func filterHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère les filtres
	filters := getFilters(r)

	// Charge toutes les données JSON
	data := loadData()

	// vérifie s'il y a des données dans les filtres sur la page HTML
	if filters["creationDate"] == "" && filters["firstAlbum"] == "" && filters["members"] == "" && filters["CitySearch"] == "" {
		fmt.Println("Pas de filtre")
		// s'il n'y a pas de filtres, on exécute les templates avec les données des fichiers Json (rendu comme mainPage)
		filterTemplate.Execute(w, data.ArtistsData)
	} else {
		// s'il y a des filtres, on appelle la fonction dataToPush et exécute le template avec filtre
		toPush := dataToPush(filters, data)

		// Rend le template avec les artistes correspondant aux filtres
		filterTemplate.Execute(w, toPush)
	}

}

// la fonction dataToPush teste les filtres sur chaque artiste et enregistre les artistes correspondants dans okFilters
func dataToPush(filters map[string]string, data Datas) []Artists {
	var okFilters Datas
	// parcourt les données de ArtistsData et compare avec tous les filtres
	for id, arti := range data.ArtistsData {
		// essaye pour chaque artiste s'il correspond aux filtres
		if testMembers, _ := strconv.Atoi(filters["members"]); len(arti.MEMBERS) == testMembers {
			// filtre suivant (Nombre de membres)
			if testDate, _ := strconv.Atoi(filters["creationDate"]); arti.CREA_DATE == testDate || testDate == 0 {
				// filtre suivant (Date de création)
				if arti.FIRST_ALBUM == filters["firstAlbum"] || filters["firstAlbum"] == "" {
					// filtre suivant (Premier Album)
					for _, City := range data.LocationsData[id].LOCATIONS {
						// filtre suivant (Lieux)
						if filters["CitySearch"] == City || filters["CitySearch"] == "" {
							// Ajoute le résultat à okFilters de tous les filtres
							okFilters.ArtistsData = append(okFilters.ArtistsData, data.ArtistsData[id])
							// affiche si l'artiste correspond
							fmt.Println(arti.NAME)
							break
						}
					}
				}
			}
		}
	}

	return okFilters.ArtistsData
}
