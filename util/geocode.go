package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bluehawk27/medcost-go/store"
)

type GeocodApi struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Bounds struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"bounds"`
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

// https://maps.googleapis.com/maps/api/geocode/json?address=1600+Amphitheatre+Parkway,
// +Mountain+View,+CA&key=YOUR_API_KEY
func GeocodeAddress(p store.Provider) (float64, float64, error) {
	key := os.Getenv("API_KEY")
	street := strings.Replace(*p.Street, " ", "+", -1)
	city := strings.Replace(*p.City, " ", "+", -1)
	state := strings.Replace(*p.State, " ", "+", -1)
	log.Info(street, city, state)
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s,%s,+%s&key=%s", street, city, state, key)
	log.Info(url)
	res, err := makeRequest(url)
	if err != nil {
		log.Error("Error geocoding the address", err)
		return 0.0, 0.0, err
	}

	log.Info(res)
	lat, long, err := extractLatLong(res)
	if err != nil {
		log.Error(err)
		return 0.0, 0.0, err
	}

	return lat, long, nil
}

func makeRequest(url string) (GeocodApi, error) {
	result := GeocodApi{}
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return result, err
	}

	json.NewDecoder(resp.Body).Decode(&result)

	return result, nil
}

func extractLatLong(body GeocodApi) (float64, float64, error) {
	if len(body.Results) < 1 {
		return 0.0, 0.0, errors.New("The Geocode api returned 0 results")
	}

	result := body.Results[0]
	lat := result.Geometry.Location.Lat
	long := result.Geometry.Location.Lng

	return lat, long, nil
}
