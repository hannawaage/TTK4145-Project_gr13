package sync

import (
	"fmt"

	"../config"
	"../driver-go/elevio"
)

// CostFunction tar inn en allElevs, id, lage ny
func CostFunction(id int, allElevs [config.NumElevs]config.Elevator, onlineIDs []int) [config.NumElevs][config.NumFloors][config.NumButtons]int {
	var allElevsMat [config.NumElevs][config.NumFloors][config.NumButtons]int

	bestElevator := allElevs[0].Id
	for elevator := 0; elevator < config.NumElevs; elevator++ {
		for floor := 0; floor < config.NumFloors; floor++ {
			for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++ {
				if allElevs[elevator].Orders[floor][button] == 1 {
					bestElevator = costCalculator(id, floor, &allElevs, onlineIDs, elevator)
					fmt.Println("bestElevator =", bestElevator)
					allElevs[elevator].Orders[floor][button] = 0
					allElevs[bestElevator].Orders[floor][button] = 2
				}
			}
		}
	}
	for elevator := 0; elevator < config.NumElevs; elevator++ {
		allElevsMat[elevator] = allElevs[elevator].Orders
	}
	return allElevsMat
}

func costCalculator(id int, floor int, allElevs *[config.NumElevs]config.Elevator, onlineIDs []int, elev int) int {
	minCost := 4*(config.NumButtons * config.NumFloors) * config.NumElevs
	bestElevator := onlineIDs[0]
	for elevator := 0; elevator < config.NumElevs; elevator++ {
		if (!Contains(onlineIDs, allElevs[elevator].Id) && (elevator != id)) {
			continue
		}
		cost := (floor - allElevs[elevator].Floor)
		if (cost == 0) && (allElevs[elevator].State != config.Moving) {
			bestElevator = elevator
			return bestElevator
		}

		if cost < 0 {
			cost = -cost
			if allElevs[elevator].Dir != elevio.MD_Down {
				cost += 3
			}
		} else if cost > 0 {
			if allElevs[elevator].Dir != elevio.MD_Up {
				cost += 3
			}
		}
		if cost == 0 && allElevs[elevator].State == config.Moving {
			cost += 5
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

func Contains(elevs []int, new int) bool {
	for _, a := range elevs {
		if a == new {
			return true
		}
	}
	return false
}

func MergeAllOrders(id int, all [config.NumElevs][config.NumFloors][config.NumButtons]int) [config.NumElevs][config.NumFloors][config.NumButtons]int {
	var merged [config.NumElevs][config.NumFloors][config.NumButtons]int
	merged[id] = all[id]
	for elev := 0; elev < config.NumElevs; elev++ {
		if elev == id {
			continue
		}
		for floor := 0; floor < config.NumFloors; floor++ {
			for btn := 0; btn < config.NumButtons; btn++ {
				if all[elev][floor][btn] > 0 && btn != config.NumButtons-1 {
					merged[id][floor][btn] = all[elev][floor][btn]
					merged[elev][floor][btn] = 0
				}
			}
		}
	}
	return merged
}

func UpdateTimeStamp(timeStamps *[config.NumFloors]int, current *[config.NumElevs][config.NumFloors][config.NumButtons]int, allElevs *[config.NumElevs]config.Elevator) {
	var numOrders int
	for floor := 0; floor < config.NumFloors; floor++ {
		numOrders = 0
		for elev := 0; elev < config.NumElevs; elev++ {
            for btn := 0; btn < config.NumButtons; btn++ {
                if (current[elev][floor][btn] > 0 ){
					numOrders ++
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

func TimeStampTimeout(timeStamps *[config.NumFloors]int) bool {
    for floor := 0; floor < config.NumFloors; floor++ {
        if timeStamps[floor] > 20 {
            return true
        }
    }
    return false
}

func FindFaultyElev(current *[config.NumElevs][config.NumFloors][config.NumButtons]int, timeStamps *[config.NumFloors]int) int {
	fmt.Println("Timestamps:")
	fmt.Println(timeStamps)
	for elev := 0; elev < config.NumElevs; elev++ {
        for floor := 0; floor < config.NumFloors; floor++ {
            for btn := 0; btn < config.NumButtons; btn++ {
                if (timeStamps[floor] > 20) && (current[elev][floor][btn] > 20) {
					return elev
        		}
            }
        }
	}
	return -1
}