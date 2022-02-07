package api

type SSEvent struct {
	ItemType string      `json:"itemType"`
	Item     interface{} `json:"item"`
	Event    string      `json:"event"`
}
