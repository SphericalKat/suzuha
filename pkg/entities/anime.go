package entities

type Anime struct {
	Url           *string                  `json:"url"`
	Type          *string                  `json:"type"`
	TrailerUrl    *string                  `json:"trailer_url"`
	Title         *string                  `json:"title"`
	TitleEnglish  *string                  `json:"title_english"`
	TitleJapanese *string                  `json:"title_japanese"`
	TitleSynonyms []string                 `json:"title_synonyms"`
	Synopsis      *string                  `json:"synopsis"`
	MalId         *int                     `json:"mal_id"`
	ImageUrl      *string                  `json:"image_url"`
	Episodes      *int                     `json:"episodes"`
	Broadcast     *string                  `json:"broadcast"`
	Duration      *string                  `json:"duration"`
	Favorites     *int                     `json:"favorites"`
	Members       *int                     `json:"members"`
	Popularity    *int                     `json:"popularity"`
	Rank          *int                     `json:"rank"`
	Score         *float64                 `json:"score"`
	ScoredBy      *int                     `json:"scored_by"`
	Rating        *string                  `json:"rating"`
	Airing        bool                     `json:"airing"`
	Aired         Aired                    `json:"aired"`
	Background    *string                  `json:"background"`
	EndingThemes  []string                 `json:"ending_themes"`
	Genres        []Genre                  `json:"genres"`
	Licensors     []Studio                 `json:"licensors"`
	OpeningThemes []string                 `json:"opening_themes"`
	Premiered     *string                  `json:"premiered"`
	Producers     []Studio                 `json:"producers"`
	Related       map[string][]RelatedItem `json:"related"`
	Source        *string                  `json:"source"`
	Status        *string                  `json:"status"`
	Studios       []Studio                 `json:"studios"`
}
