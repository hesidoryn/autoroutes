package tracer

import (
	"context"
	"time"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmapi"
)

// Tracer is used to plot a route between two nodes
type Tracer struct {
	waysCache map[osm.WayID]*osm.Way
}

// New returns Tracer
func New() *Tracer {
	return &Tracer{
		waysCache: map[osm.WayID]*osm.Way{},
	}
}

// Trace returns all ways between two nodes
func (t *Tracer) Trace(from, to int64) (*osm.Way, error) {
	time.Sleep(1 * time.Second)
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()
	waysFrom, err := osmapi.NodeWays(ctx, osm.NodeID(from))
	if err != nil {
		return nil, err
	}
	waysTo, err := osmapi.NodeWays(ctx, osm.NodeID(to))
	if err != nil {
		return nil, err
	}

	for _, idFrom := range waysFrom.IDs() {
		for _, idTo := range waysTo.IDs() {
			if idFrom == idTo {
				way, ok := t.waysCache[idFrom]
				if ok {
					return way, nil
				}

				way, err = osmapi.Way(ctx, idFrom)
				if err != nil {
					return nil, err
				}

				if way.Tags.Find("highway") == "service" &&
					way.Tags.Find("living_street") == "yes" {
					return nil, nil
				}

				t.waysCache[idFrom] = way
				return way, nil
			}
		}
	}

	return nil, nil
}
