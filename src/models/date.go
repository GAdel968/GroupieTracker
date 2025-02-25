package models

type Date struct {
	ID    int                 `json:"id"`
	Dates map[string][]string `json:"dates"`
}
