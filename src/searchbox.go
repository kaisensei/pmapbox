package mapbox

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

type Location struct {
	MapboxID     string
	Name         string
	FullCityName string
	Lat          float64
	Lng          float64
}

// https://docs.mapbox.com/data/boundaries/guides/boundaries-v4-migration-guide/
func RetrieveByMapboxID(mapboxID, sessionToken string) ([]*Location, error) {
	locations := make([]*Location, 3)
	accessToken := viper.Get("MAPBOX_TOKEN").(string)
	if accessToken == "" {
		return nil, fmt.Errorf("MAPBOX_ACCESS_TOKEN is missing")
	}
	resp, err := http.Get(getQueryUrl(mapboxID, accessToken, sessionToken))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	var parsedResp SearchboxReverseResponse
	if err := json.Unmarshal(body, &parsedResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}
	feature := parsedResp.Features[0]
	locations[0] = getLocationFromFeature(feature, mapboxID)
	for key, value := range feature.Properties.Context {
		if key == "street" {
			locations[1], err = retrieveByMapboxIDUtil(value.ID, accessToken, sessionToken)
			if err != nil {
				return nil, err
			}
		}
		if key == "place" {
			locations[2], err = retrieveByMapboxIDUtil(value.ID, accessToken, sessionToken)
			if err != nil {
				return nil, err
			}
		}
	}
	return locations, nil
}

func GetMapboxSuggestions(queryStr, sessionToken string) ([]*MapboxSuggestion, error) {
	accessToken := viper.Get("MAPBOX_TOKEN").(string)
	if accessToken == "" {
		return nil, fmt.Errorf("MAPBOX_ACCESS_TOKEN is missing")
	}
	queryUrl := fmt.Sprintf("https://api.mapbox.com/search/searchbox/v1/suggest?q=%s&session_token=%s&access_token=%s", url.QueryEscape(queryStr), sessionToken, accessToken)
	resp, err := http.Get(queryUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	var suggestions MapboxSuggestions
	if err := json.Unmarshal(body, &suggestions); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}
	if len(suggestions.Suggestions) > 6 {
		return suggestions.Suggestions[:6], nil
	}
	return suggestions.Suggestions, nil
}

func retrieveByMapboxIDUtil(mapboxID, accessToken, sessionToken string) (*Location, error) {
	resp, err := http.Get(getQueryUrl(mapboxID, accessToken, sessionToken))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	var parsedResp SearchboxReverseResponse
	if err := json.Unmarshal(body, &parsedResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}
	return getLocationFromFeature(parsedResp.Features[0], mapboxID), nil
}

func getLocationFromFeature(feature *SearchboxFeature, mapboxID string) *Location {
	coordindates := feature.Geometry.Coordinates
	properties := feature.Properties
	return &Location{
		MapboxID:     mapboxID,
		Name:         properties.Name,
		FullCityName: properties.FullCityName,
		Lat:          coordindates[1],
		Lng:          coordindates[0],
	}
}

func getQueryUrl(mapboxID, accessToken, sessionToken string) string {
	return fmt.Sprintf("https://api.mapbox.com/search/searchbox/v1/retrieve/%s?access_token=%s&session_token=%s", mapboxID, accessToken, sessionToken)
}
