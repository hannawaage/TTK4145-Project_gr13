package sync

import (
	"../config"
	"../driver-go/elevio"
)

// CostFunction tar inn en allElevs, id, lage ny
func CostFunction(id int, allElevs [config.NumElevs]config.Elevator, onlineIDs []int) [config.NumElevs][config.NumFloors][config.NumButtons]bool {
	var allElevsMat [config.NumElevs][config.NumFloors][config.NumButtons]bool
	bestElevator := allElevs[0].Id
	for elevator := 0; elevator < config.NumElevs; elevator++ {
		for floor := 0; floor < config.NumFloors; floor++ {
			for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++ {
				if allElevs[elevator].Orders[floor][button] == true {
					order := elevio.ButtonEvent{
						Floor:  floor,
						Button: button,
					}
					bestElevator = costCalculator(id, order, &allElevs, onlineIDs)
					allElevs[elevator].Orders[order.Floor][order.Button] = false
					allElevs[bestElevator].Orders[order.Floor][order.Button] = true
				}
			}
		}
	}
	for elevator := 0; elevator < config.NumElevs; elevator++ {
		allElevsMat[elevator] = allElevs[elevator].Orders
	}
	return allElevsMat
}

func costCalculator(id int, order elevio.ButtonEvent, allElevs *[config.NumElevs]config.Elevator, onlineIDs []int) int {
	minCost := (config.NumButtons * config.NumFloors) * config.NumElevs
	bestElevator := id
	for elevator := 0; elevator < config.NumElevs; elevator++ {
		if !contains(onlineIDs, allElevs[elevator].Id) && (elevator != id) {
			continue
		}
		cost := order.Floor - allElevs[elevator].Floor
		if (cost == 0) && (allElevs[elevator].State != config.Moving) {
			bestElevator = elevator
			return bestElevator
		}

		if cost < 0 {
			cost = -cost
			if allElevs[elevator].Dir == elevio.MD_Up {
				cost += 3
			}
		} else if cost > 0 {
			if allElevs[elevator].Dir == elevio.MD_Down {
				cost += 3
			}
		}
		if cost == 0 && allElevs[elevator].State == config.Moving {
			cost += 10
		}

		if allElevs[elevator].State == config.DoorOpen {
			cost++
		}

		if cost < minCost {
			minCost = cost
			bestElevator = elevator
		}
	}
	return bestElevator
}

func contains(elevs []int, new int) bool {
	for _, a := range elevs {
		if a == new {
			return true
		}
	}
	return false
}

func mergeAllOrders(id int, all [config.NumElevs][config.NumFloors][config.NumButtons]bool) [config.NumElevs][config.NumFloors][config.NumButtons]bool {
	var merged [config.NumElevs][config.NumFloors][config.NumButtons]bool
	merged[id] = all[id]
	for elev := 0; elev < config.NumElevs; elev++ {
		if elev == id {
			continue
		}
		for floor := 0; floor < config.NumFloors; floor++ {
			for btn := 0; btn < config.NumButtons; btn++ {
				if all[elev][floor][btn] && btn != config.NumButtons-1 {
					merged[id][floor][btn] = true
					merged[elev][floor][btn] = false
				}
			}
		}
	}
	return merged
}

func newCabOrdersOnly(id int, current *[config.NumElevs][config.NumFloors][config.NumButtons]bool, updated *[config.NumElevs][config.NumFloors][config.NumButtons]bool) bool {
	var newCab bool
	for floor := 0; floor < config.NumFloors; floor++ {
		for btn := 0; btn < config.NumButtons-1; btn++ {
			if current[id][floor][btn] != updated[id][floor][btn] {
				return false
			}
		}
		if current[id][floor][2] != updated[id][floor][2] {
			newCab = true
		}
	}
	return newCab
}

func mergeLocalOrders(id int, local *[config.NumFloors][config.NumButtons]bool, incomming [config.NumElevs][config.NumFloors][config.NumButtons]bool) [config.NumElevs][config.NumFloors][config.NumButtons]bool {
	merged := incomming
	for floor := 0; floor < config.NumFloors; floor++ {
		for btn := 0; btn < config.NumButtons-1; btn++ {
			if local[floor][btn] {
				merged[id][floor][btn] = true
			}
		}
	}
	return merged
}
