package main

import "github.com/paulmach/osm"

func removeDuplicates(elements []*osm.Way) []*osm.Way {
	// Use map to record duplicates as we find them.
	encountered := map[*osm.Way]bool{}
	result := []*osm.Way{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}
