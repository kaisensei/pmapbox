package mapbox

type MapboxSuggestion struct {
	MapboxID    string `json:"mapbox_id"`
	FullAddress string `json:"full_address"`
}

type MapboxSuggestions struct {
	Suggestions []*MapboxSuggestion `json:"suggestions,omitempty"`
}
