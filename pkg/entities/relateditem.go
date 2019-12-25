package entities

type RelatedItem struct {
	MalId int    `json:"mal_id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	Url   string `json:"url"`
}
