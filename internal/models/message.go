package models

type Update struct {
	Device  string
	Element string `json:"elem"`
	Data    string `json:"data"`
}
