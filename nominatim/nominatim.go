package nominatim

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	apiURL  = "https://nominatim.openstreetmap.org/%v"
	reverse = "reverse?format=jsonv2&lat=%v&lon=%v&zoom=18"
	search  = "search/%v?format=json&addressdetails=1&limit=1&zoom=18"
)

// Address describes some OpenStreetMap object
type Address struct {
	PlaceID     string `json:"place_id"`
	OsmType     string `json:"osm_type"`
	OsmID       string `json:"osm_id"`
	ID          int64
	Latitude    string `json:"lat"`
	Lat         float64
	Longitude   string `json:"lon"`
	Lon         float64
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

	endpoint := fmt.Sprintf(reverse, lat, lon)
	url := fmt.Sprintf(apiURL, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := g.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(address)
	if err != nil {
		return nil, err
	}

	address.ID, err = strconv.ParseInt(address.OsmID, 10, 64)
	fmt.Println(address.ID)
	if err != nil {
		return nil, err
	}
	address.Lat, err = strconv.ParseFloat(address.Latitude, 64)
	if err != nil {
		return nil, err
	}
	address.Lon, err = strconv.ParseFloat(address.Longitude, 64)
	if err != nil {
		return nil, err
	}
	return address, nil
	//
	// // address.Latitude, _ = strconv.ParseFloat(address.Lat, 64)
	// // address.Longitude, _ = strconv.ParseFloat(address.Lon, 64)
	// return address, nil
}

// Search is used to get object address by query
func (g *Geocoder) Search(query string) (*Address, error) {
	time.Sleep(1 * time.Second)
	addresses := []*Address{}

	endpoint := fmt.Sprintf(search, query)
	url := fmt.Sprintf(apiURL, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := g.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&addresses)
	if err != nil {
		return nil, err
	}

	address := addresses[0]
	address.ID, err = strconv.ParseInt(address.OsmID, 10, 64)
	if err != nil {
		return nil, err
	}
	address.Lat, err = strconv.ParseFloat(address.Latitude, 64)
	if err != nil {
		return nil, err
	}
	address.Lon, err = strconv.ParseFloat(address.Longitude, 64)
	if err != nil {
		return nil, err
	}
	return address, nil
}
