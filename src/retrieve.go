package mapbox

type SearchboxReverseResponse struct {
	Features []*SearchboxFeature `json:"features"`
}

type MapboxContext struct {
	// Always present
	ID   string `json:"id"`
	Name string `json:"name"`

	// Region fields
	RegionCode     string `json:"region_code,omitempty"`
	RegionCodeFull string `json:"region_code_full,omitempty"`

	// Address fields
	AddressNumber string `json:"address_number,omitempty"`
	StreetName    string `json:"street_name,omitempty"`

	// Country fields
	CountryCode       string `json:"country_code,omitempty"`
	CountryCodeAlpha3 string `json:"country_code_alpha_3,omitempty"`
}

type SearchboxFeature struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Geometry   *Geometry   `json:"geometry"`
	Properties *Properties `json:"properties,omitempty"`
}

type Properties struct {
	MapboxID     string                   `json:"mapbox_id"`
	FeatureType  string                   `json:"feature_type"`
	Name         string                   `json:"name"`
	FullCityName string                   `json:"full_city_name"`
	Context      map[string]MapboxContext `json:"context,omitempty"`
	BoundingBox  []float64                `json:"bbox,omitempty"`
}
