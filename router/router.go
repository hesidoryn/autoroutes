package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	osrm "github.com/gojuno/go.osrm"
	geo "github.com/paulmach/go.geo"
)

const api = "https://router.project-osrm.org/route/v1/driving/%v,%v;%v,%v?overview=false&annotations=true"

// Router is used to create route
type Router struct {
	c *osrm.OSRM
}

// New create new Router
func New() *Router {
	return &Router{
		c: osrm.NewFromURL("https://router.project-osrm.org"),
	}
}

// Route returns route between some points
func (r *Router) Route(points geo.PointSet) (*osrm.RouteResponse, error) {
	// ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancelFn()

	url := fmt.Sprintf(api,
		points.First().Lng(),
		points.First().Lat(),
		points.Last().Lng(),
		points.Last().Lat(),
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := &osrm.RouteResponse{}
	return response, json.NewDecoder(resp.Body).Decode(response)

	// return r.c.Route(ctx, osrm.RouteRequest{
	// 	Profile:     "driving",
	// 	GeoPath:     *osrm.NewGeoPathFromPointSet(points),
	// 	Steps:       osrm.StepsFalse,
	// 	Annotations: osrm.AnnotationsTrue,
	// 	Overview:    osrm.OverviewFalse,
	// 	Geometries:  osrm.GeometriesPolyline6,
	// })
}
