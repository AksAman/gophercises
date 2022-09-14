package models

type URLSet struct {
	URLs []URL `json:"url"`
}

type URL struct {
	Loc        string  `json:"loc"`
	Lastmod    *string `json:"lastmod,omitempty"`
	Changefreq *string `json:"changefreq,omitempty"`
	Priority   float64 `json:"priority"`
}
