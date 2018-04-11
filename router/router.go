package router

import (
	"context"
	"time"

	osrm "github.com/gojuno/go.osrm"
	geo "github.com/paulmach/go.geo"
)

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
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	return r.c.Route(ctx, osrm.RouteRequest{
		Profile:     "driving",
		GeoPath:     *osrm.NewGeoPathFromPointSet(points),
		Steps:       osrm.StepsFalse,
		Annotations: osrm.AnnotationsTrue,
		Overview:    osrm.OverviewFalse,
		Geometries:  osrm.GeometriesPolyline6,
	})
}
