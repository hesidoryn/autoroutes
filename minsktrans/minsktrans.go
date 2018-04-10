package minsktrans

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	apiURL    = "https://minsktrans-api-ggodcuvsbr.now.sh"
	getRoutes = "routes"
	getStops  = "stops"
)

// Route describes public transport route
type Route struct {
	ID              string   `json:"id"`
	RouteNum        string   `json:"routeNum"`
	Operator        string   `json:"operator"`
	ValidityPeriods string   `json:"validityPeriods"`
	RouteType       string   `json:"routeType"`
	RouteName       string   `json:"routeName"`
	Weekdays        string   `json:"weekdays"`
	Stops           []string `json:"stops"`
	Datestart       string   `json:"datestart"`
}

// Stop describes public transport stop
type Stop struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Longitude float64  `json:"lng"`
	Latitude  float64  `json:"lat"`
	Stops     []string `json:"stops"`
}

// Client is minsktrans api client
type Client struct {
	c      *http.Client
	routes []*Route
	stops  []*Stop
}

// New creates minsktrans api client
func New() *Client {
	return &Client{
		c: &http.Client{},
	}
}

// LoadRoutes loads Minsk's public transport routes
func (cl *Client) LoadRoutes() error {
	routes := []*Route{}
	err := cl.makeRequest(getRoutes, &routes)
	if err != nil {
		return err
	}
	cl.routes = routes
	return err
}

// LoadStops loads Minsk's public transport stops
func (cl *Client) LoadStops() error {
	stops := []*Stop{}
	err := cl.makeRequest(getStops, &stops)
	if err != nil {
		return err
	}
	cl.stops = stops
	return err
}

// GetRoutes returns all routes
func (cl *Client) GetRoutes() []*Route {
	return cl.routes
}

// GetRoute returns route by id
func (cl *Client) GetRoute(id string) *Route {
	for _, r := range cl.routes {
		if r.ID == id {
			return r
		}
	}
	return &Route{}
}

// GetStop returns route by id
func (cl *Client) GetStop(id string) *Stop {
	for _, s := range cl.stops {
		if s.ID == id {
			return s
		}
	}
	return &Stop{}
}

func (cl *Client) makeRequest(endpoint string, output interface{}) error {
	url := fmt.Sprintf("%s/%s", apiURL, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := cl.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body) // todo: remove _
		return fmt.Errorf("Response body: %v", string(body))
	}

	return json.NewDecoder(resp.Body).Decode(output)
}
