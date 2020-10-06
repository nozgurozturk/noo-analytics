package entities

type dateRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type date struct {
	IsInclude bool      `json:"isInclude"`
	Range     dateRange `json:"range"`
}

type AnalyticsActionRequest struct {
	Action Action `json:"action"`
	Year   date   `json:"year"`
	Month  date   `json:"month"`
	Day    date   `json:"day"`
	Hour   date   `json:"hour"`
}

type AnalyticsActionResponse struct {
	Year           int      `json:"y" bson:"y, omitempty"`
	Month          int      `json:"m" bson:"m, omitempty"`
	Day            int      `json:"d" bson:"d, omitempty"`
	Hour           int      `json:"h" bson:"h, omitempty"`
	Visitors       []string `json:"v" bson:"v, omitempty"`
	UniqueVisitors []string `json:"uv" bson:"uv, omitempty"`
}
