package models

type Location struct {
	ID        int                 `json:"id"`
	Locations map[string][]string `json:"locations"`
	DatesURL  string              `json:"dates"`
}
