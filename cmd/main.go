package main

import (
	"fmt"
	"log"

	"github.com/hesidoryn/autoroutes/minsktrans"
	"github.com/hesidoryn/autoroutes/nominatim"
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

	geocoder := nominatim.New()

	routes := mt.GetRoutes()
	fmt.Println("Route count: ", len(routes))

	fullCount := 0
	for _, r := range routes {
		stopCount := 0
		for _, id := range r.Stops {
			stop := mt.GetStop(id)
			// fmt.Println(stop.Latitude, stop.Longitude)
			address, _ := geocoder.ReverseGeocode(stop.Latitude, stop.Longitude)
			// fmt.Println(address)
			if address.Type == "bus_stop" ||
				address.Type == "bus_station" {
				stopCount++
			}
		}
		if stopCount == len(r.Stops) {
			fullCount++
		}
		fmt.Printf("Route name: %v, expected stops: %v, actual stops: %v\n", r.RouteName, len(r.Stops), stopCount)
	}

	fmt.Println(fullCount)

	// geocoder := nominatim.New()
	// address, _ := geocoder.ReverseGeocode(stops[101].Latitude, stops[101].Longitude)
	// if address != nil {
	// 	fmt.Printf("Detailed address: %#v\n", address)
	// } else {
	// 	fmt.Println("got <nil> address")
	// }
}
