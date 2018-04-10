package nominatim

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const apiURL = "https://nominatim.openstreetmap.org/reverse?format=jsonv2&lat=%v&lon=%v&zoom=18"

// Address describes some OpenStreetMap object
type Address struct {
	PlaceID  string `json:"place_id"`
	OsmType  string `json:"osm_type"`
	OsmID    string `json:"osm_id"`
	Latitude string `json:"lat"`
	// Latitude    float64
	Longitude string `json:"lon"`
	// Longitude   float64
	Category    string `json:"category"`
	Type        string `json:"type"`
	AddressType string `json:"addresstype"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Address     struct {
		BusStop      string `json:"bus_stop"`
		Road         string `json:"road"`
		Suburb       string `json:"suburb"`
		CityDistrict string `json:"city_district"`
		City         string `json:"city"`
		PostCode     string `json:"postcode"`
		Country      string `json:"country"`
		CountryCode  string `json:"country_code"`
	} `json:"address"`
}

// Geocoder is OpenStreetMap geocoder client
type Geocoder struct {
	c *http.Client
}

// New creates OpenStreetMap geocoder
func New() *Geocoder {
	return &Geocoder{
		c: &http.Client{},
	}
}

// ReverseGeocode generates an address from a latitude and longitude
func (g *Geocoder) ReverseGeocode(lat, lon float64) (*Address, error) {
	time.Sleep(1 * time.Second)
	address := &Address{}

	url := fmt.Sprintf(apiURL, lat, lon)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return address, err
	}

	resp, err := g.c.Do(req)
	if err != nil {
		return address, err
	}
	defer resp.Body.Close()

	return address, json.NewDecoder(resp.Body).Decode(address)
	//
	// // address.Latitude, _ = strconv.ParseFloat(address.Lat, 64)
	// // address.Longitude, _ = strconv.ParseFloat(address.Lon, 64)
	// return address, nil
}
