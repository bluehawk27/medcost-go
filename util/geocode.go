package util

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bluehawk27/medcost-go/store"
)

// https://maps.googleapis.com/maps/api/geocode/json?address=1600+Amphitheatre+Parkway,
// +Mountain+View,+CA&key=YOUR_API_KEY
func GeocodeAddress(p store.Provider) (float64, float64, error) {
	// p.Street
	// p.City
	// p.State
	// p.Zipcode
	fmt.Println(p.Street)
	fmt.Println(p.City)
	fmt.Println(p.State)
	key := os.Getenv("API_KEY")
	fmt.Println(key)
	// url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address%s, %s", e.When, e.What)
	req, err := http.NewRequest("GET", "https://yahoo.com", nil)
	if err != nil {
		return 0, 0, err
	}
	fmt.Println(req)
	lat := 00.00
	long := 00.00
	return lat, long, nil
}
