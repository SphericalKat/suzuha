package entities

type IndexInfo struct {
	Author           string
	Telegram         string
	Version          string
	SuzuhaGo       string
	Website          string
	Docs             string
	GitHub           string
	ProductionApiUrl string `json:"PRODUCTION_API_URL"`
	StatusUrl        string `json:"STATUS_URL"`
}
