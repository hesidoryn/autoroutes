package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hesidoryn/autoroutes/minsktrans"
	"github.com/hesidoryn/autoroutes/nominatim"
	"github.com/hesidoryn/autoroutes/router"
	"github.com/hesidoryn/autoroutes/tracer"
	geo "github.com/paulmach/go.geo"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmapi"
)

func main() {
	mt := minsktrans.New()
	err := mt.LoadRoutes()
	if err != nil {
		log.Fatal(err)
	}

	err = mt.LoadStops()
	if err != nil {
		log.Fatal(err)
	}

	relation107 := &osm.Relation{}

	geocoder := nominatim.New()
	route107 := mt.GetRoute("208392")
	stops107 := []*nominatim.Address{}
	for _, id := range route107.Stops {
		stop := mt.GetStop(id)
		address, err := geocoder.ReverseGeocode(stop.Latitude, stop.Longitude)
		if err != nil {
			log.Fatal(err)
		}
		if address.Type == "bus_stop" ||
			address.Type == "bus_station" {
			stops107 = append(stops107, address)
			continue
		}
		query := fmt.Sprintf("%v %v, %v", stop.Latitude, stop.Longitude, stop.Name)
		address, _ = geocoder.Search(query)
		stops107 = append(stops107, address)
	}

	for i, stop := range stops107 {
		member := osm.Member{
			Type: osm.TypeNode,
			Ref:  stop.ID,
			Role: "platform",
		}
		if i == 0 {
			member.Role += "_entry_only"
		}
		if i == 0 {
			member.Role += "_exit_only"
		}
		relation107.Members = append(relation107.Members, member)
	}

	router := router.New()
	tracer := tracer.New()
	ways107 := []*osm.Way{}
	for i := 0; i < len(stops107)-1; i++ {
		points := geo.PointSet{
			{stops107[i].Lon, stops107[i].Lat},
			{stops107[i+1].Lon, stops107[i+1].Lat},
		}
		fmt.Println(points)
		route, err := router.Route(points)
		if err != nil {
			log.Fatal(err)
		}

		keyNodeIDs := route.Routes[0].Legs[0].Annotation.Nodes
		for i := 0; i < len(keyNodeIDs)-1; i++ {
			way, err := tracer.Trace(keyNodeIDs[i], keyNodeIDs[i+1])
			if err != nil {
				log.Fatal(err)
			}
			if way == nil {
				continue
			}
			ways107 = append(ways107, way)
		}
		last := osm.NodeID(keyNodeIDs[len(keyNodeIDs)-3])
		lastNode, err := osmapi.Node(context.Background(), last)
		if err != nil {
			log.Fatalf("Check last node err: %v", err)
		}
		if lastNode.Tags.Find("public_transport") != "" {
			fmt.Println(lastNode.Tags.Find("public_transport"))
			continue
		}
		last = osm.NodeID(keyNodeIDs[len(keyNodeIDs)-2])
		lastNode, err = osmapi.Node(context.Background(), last)
		if err != nil {
			log.Fatalf("Check last node err: %v", err)
		}
		if lastNode.Tags.Find("public_transport") != "" {
			fmt.Println(lastNode.Tags.Find("public_transport"))
			continue
		}
		last = osm.NodeID(keyNodeIDs[len(keyNodeIDs)-1])
		lastNode, err = osmapi.Node(context.Background(), last)
		if err != nil {
			log.Fatalf("Check last node err: %v", err)
		}
		fmt.Println(lastNode.Tags.Find("public_transport"))
	}
	fmt.Println("---------------------------")
	for _, way := range removeDuplicates(ways107) {
		fmt.Println(way.ID)
	}

	// geocoder := nominatim.New()
	//
	// routes := mt.GetRoutes()
	// fmt.Println("Route count: ", len(routes))
	//
	// fullCount := 0
	// for _, r := range routes {
	// 	stopCount := 0
	// 	for _, id := range r.Stops {
	// 		stop := mt.GetStop(id)
	// 		fmt.Println(stop.Latitude, stop.Longitude)
	// 		address, _ := geocoder.ReverseGeocode(stop.Latitude, stop.Longitude)
	// 		fmt.Println(address)
	// 		if address.Type == "bus_stop" ||
	// 			address.Type == "bus_station" {
	// 			stopCount++
	// 		}
	// 	}
	// 	if stopCount == len(r.Stops) {
	// 		fullCount++
	// 	}
	// 	fmt.Printf("Route name: %v, expected stops: %v, actual stops: %v\n", r.RouteName, len(r.Stops), stopCount)
	// }
	//
	// fmt.Println(fullCount)

	// geocoder := nominatim.New()
	// address, _ := geocoder.ReverseGeocode(stops[101].Latitude, stops[101].Longitude)
	// if address != nil {
	// 	fmt.Printf("Detailed address: %#v\n", address)
	// } else {
	// 	fmt.Println("got <nil> address")
	// }
}
