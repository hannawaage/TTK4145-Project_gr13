package sync

import (
	"../config"
	"../driver-go/elevio"
	"fmt"
)

const (
	Moving   = config.Moving
	DoorOpen = config.DoorOpen
)

func CostFunction(id int, allElevs [NumElevs]config.Elevator, onlineIDs []int) [NumElevs][NumFloors][NumButtons]int {
	var allElevsMat [NumElevs][NumFloors][NumButtons]int

	bestElevator := allElevs[0].Id
	for elevator := 0; elevator < NumElevs; elevator++ {
		for floor := 0; floor < NumFloors; floor++ {
			for btn := elevio.BT_HallUp; btn < elevio.BT_Cab; btn++ {
				if allElevs[elevator].Orders[floor][btn] == 1 {
					bestElevator = costCalculator(id, floor, &allElevs, onlineIDs)
					allElevs[elevator].Orders[floor][btn] = 0
					allElevs[bestElevator].Orders[floor][btn] = 2
				}
			}
		}
	}
	for elevator := 0; elevator < NumElevs; elevator++ {
		allElevsMat[elevator] = allElevs[elevator].Orders
	}
	return allElevsMat
}

func costCalculator(id int, floor int, allElevs *[NumElevs]config.Elevator, onlineIDs []int) int {
	minCost := 4*(NumButtons * NumFloors) * NumElevs
	bestElevator := onlineIDs[0]
	for elevator := 0; elevator < NumElevs; elevator++ {
		if !Contains(onlineIDs, allElevs[elevator].Id) && (elevator != id) {
			continue
		}
		cost := 2*(floor - allElevs[elevator].Floor)
		if (cost == 0) && (allElevs[elevator].State != Moving) {
			bestElevator = elevator
			return bestElevator
		}

		if cost < 0 {
			cost = -cost
			if allElevs[elevator].Dir != elevio.MD_Down {
				cost += 1
			}
			if allElevs[elevator].Dir == elevio.MD_Up {
				cost += 1
			}
		} else if cost > 0 {
			if allElevs[elevator].Dir != elevio.MD_Up {
				cost += 1
			}
		  if allElevs[elevator].Dir == elevio.MD_Down {
				cost += 1
			}
		}
		if cost == 0 && allElevs[elevator].State == Moving {
			cost += 3
		}

		fmt.Println("cost = ",cost)
		if allElevs[elevator].State == config.DoorOpen {
			cost += 1
		}

		if cost < minCost {
			minCost = cost
			bestElevator = elevator
		}
	}
	return bestElevator
}

func Contains(elevs []int, new int) bool {
	for _, elev := range elevs {
		if elev == new {
			return true
		}
	}
	return false
}

func MergeAllOrders(id int, all [NumElevs][NumFloors][NumButtons]int) [NumElevs][NumFloors][NumButtons]int {
	var merged [NumElevs][NumFloors][NumButtons]int
	merged[id] = all[id]
	for elev := 0; elev < NumElevs; elev++ {
		if elev == id {
			continue
		}
		for floor := 0; floor < NumFloors; floor++ {
			for btn := 0; btn < NumButtons; btn++ {
				if all[elev][floor][btn] > 0 && btn != NumButtons-1 {
					merged[id][floor][btn] = all[elev][floor][btn]
					merged[elev][floor][btn] = 0
				}
			}
		}
	}
	return merged
}

func UpdateTimeStamp(timeStamps *[NumFloors]int, current *[NumElevs][NumFloors][NumButtons]int, allElevs *[NumElevs]config.Elevator) {
	var numOrders int
	for floor := 0; floor < NumFloors; floor++ {
		numOrders = 0
		for elev := 0; elev < NumElevs; elev++ {
			for btn := 0; btn < NumButtons; btn++ {
				if current[elev][floor][btn] > 0 {
					numOrders++
					timeStamps[floor]++
				}
			}
		}
		if numOrders == 0 {
			if timeStamps[floor] != 0 {
				timeStamps[floor] = 0
			}
		}
	}
}

func OrderTimeout(timeStamps *[NumFloors]int) bool {
	for floor := 0; floor < NumFloors; floor++ {
		if timeStamps[floor] > 120 {
			return true
		}
	}
	return false
}

func FindFaultyElev(current *[NumElevs][NumFloors][NumButtons]int, timeStamps *[NumFloors]int) int {
	for elev := 0; elev < NumElevs; elev++ {
		for floor := 0; floor < NumFloors; floor++ {
			for btn := 0; btn < NumButtons; btn++ {
				if (timeStamps[floor] > 120) && (current[elev][floor][btn] > 0) {
					return elev
				}
			}
		}
	}
	return -1
}
