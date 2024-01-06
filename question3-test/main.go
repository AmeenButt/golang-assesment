package main

import (
	"fmt"
	"sort"
)

// Seat struct to hold the seat data
type Seat struct {
	id      int
	student string
}

func sortByID(seats []Seat) {
	sort.Slice(seats, func(i, j int) bool {
		return seats[i].id < seats[j].id
	})
}

// swapSeats function to swap every two consecutive seats
func swapSeats(seats []Seat) []Seat {
	for i := 0; i < len(seats); i += 2 {
		if i+1 == len(seats) && len(seats)%2 != 0 {

		} else {
			seats[i].id, seats[i+1].id = seats[i+1].id, seats[i].id
		}

	}
	sortByID(seats)
	return seats
}

func main() {
	// Hardcoded input data
	seats := []Seat{
		{1, "Abbot"},
		{2, "Doris"},
		{3, "Emerson"},
		{4, "Green"},
		{5, "Jeames"},
		{6, "Ameen"},
	}

	// Perform the swap operation
	swappedSeats := swapSeats(seats)

	// Displaying the output
	for _, seat := range swappedSeats {
		fmt.Printf("id: %d, student: %s\n", seat.id, seat.student)
	}
}
