package sync

import (
	"../config"
	"../driver-go/elevio"
)

const (
	Undefined = config.Undefined
	Idle      = config.Idle
	Moving    = config.Moving
	DoorOpen  = config.DoorOpen
)

func CostFunction(id int, allElevs [NumElevs]config.Elevator, onlineIDs []int) [NumElevs][NumFloors][NumButtons]bool {
	var allElevsMat [NumElevs][NumFloors][NumButtons]bool
	bestElevator := allElevs[0].Id
	for elevator := 0; elevator < NumElevs; elevator++ {
		for floor := 0; floor < NumFloors; floor++ {
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
	for elevator := 0; elevator < NumElevs; elevator++ {
		allElevsMat[elevator] = allElevs[elevator].Orders
	}
	return allElevsMat
}

func costCalculator(id int, order elevio.ButtonEvent, allElevs *[NumElevs]config.Elevator, onlineIDs []int) int {
	minCost := (NumButtons * NumFloors) * NumElevs
	bestElevator := onlineIDs[0]
	for elevator := 0; elevator < NumElevs; elevator++ {
		if (!contains(onlineIDs, allElevs[elevator].Id) && (elevator != id)) || (allElevs[elevator].State == Undefined) {
			continue
		}
		cost := order.Floor - allElevs[elevator].Floor
		if (cost == 0) && (allElevs[elevator].State != Moving) {
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
		if cost == 0 && allElevs[elevator].State == Moving {
			cost += 4
		}

		if allElevs[elevator].State == DoorOpen {
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

func mergeAllOrders(id int, all [NumElevs][NumFloors][NumButtons]bool) [NumElevs][NumFloors][NumButtons]bool {
	var merged [NumElevs][NumFloors][NumButtons]bool
	merged[id] = all[id]
	for elev := 0; elev < NumElevs; elev++ {
		if elev == id {
			continue
		}
		for floor := 0; floor < NumFloors; floor++ {
			for btn := 0; btn < NumButtons; btn++ {
				if all[elev][floor][btn] && btn != NumButtons-1 {
					merged[id][floor][btn] = true
					merged[elev][floor][btn] = false
				}
			}
		}
	}
	return merged
}

func updateTimeStamp(timeStamps *[NumFloors]int, current *[NumElevs][NumFloors][NumButtons]bool, updated *[NumElevs][NumFloors][NumButtons]bool) {
	for elev := 0; elev < NumElevs; elev++ {
		for floor := 0; floor < NumFloors; floor++ {
			for btn := 0; btn < NumButtons; btn++ {
				if updated[elev][floor][btn] {
					timeStamps[floor]++
				} else if !updated[elev][floor][btn] && current[elev][floor][btn] {
					timeStamps[floor] = 0
				}
			}
		}
	}
}

func TimeStampTimeout(timeStamps *[NumFloors]int) bool {
	for floor := 0; floor < NumFloors; floor++ {
		if timeStamps[floor] > 10 {
			return true
		}
	}
	return false
}
