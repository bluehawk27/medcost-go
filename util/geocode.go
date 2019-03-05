package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
// assign to provider and return that ?
func GeocodeAddress(p store.Provider) (float64, float64, error) {
	key := os.Getenv("API_KEY")

	street := strings.Replace(*p.Street, " ", "+", -1)
	city := strings.Replace(*p.City, " ", "+", -1)
	state := strings.Replace(*p.State, " ", "+", -1)
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s,%s,+%s&key=%s", street, city, state, key)

	res, err := MakeRequest(url)
	if err != nil {
		log.Error("Error geocoding the address", err)
		return 0.0, 0.0, err
	}

	result := GeocodApi{}
	json.NewDecoder(res).Decode(&result)
	lat, long, err := extractLatLong(result)
	if err != nil {
		log.Error(err)
		return 0.0, 0.0, err
	}

	return lat, long, nil
}

func MakeRequest(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return resp.Body, nil
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
