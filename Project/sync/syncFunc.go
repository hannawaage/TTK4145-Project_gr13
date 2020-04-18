package sync

import (
	"fmt"

	"../config"
	"../driver-go/elevio"
)

// CostFunction tar inn en allElevs, id, lage ny
func CostFunction(id int, allElevs [config.NumElevs]config.Elevator, onlineIDs []int) [config.NumElevs][config.NumFloors][config.NumButtons]bool {
	/*
		Hvis heisen ordren ligger hos er i bevegelse, skal den ikke forandres
	*/

	var allElevsMat [config.NumElevs][config.NumFloors][config.NumButtons]bool
	for elevator := 0; elevator < config.NumElevs; elevator++ {
		allElevsMat[elevator] = allElevs[elevator].Orders
	}

	prevSum := sumOrders(allElevsMat)
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
					fmt.Println("BestElev", bestElevator)
					allElevs[elevator].Orders[order.Floor][order.Button] = false
					allElevs[bestElevator].Orders[order.Floor][order.Button] = true
				}
			}
		}
	}
	for elevator := 0; elevator < config.NumElevs; elevator++ {
		allElevsMat[elevator] = allElevs[elevator].Orders
	}
	afterSum := sumOrders(allElevsMat)

	if afterSum < prevSum {
		fmt.Println("Prevsum", prevSum)
		fmt.Println("Aftersum", afterSum)
	}
	return allElevsMat
}

func costCalculator(id int, order elevio.ButtonEvent, allElevs *[config.NumElevs]config.Elevator, onlineIDs []int) int {
	minCost := (config.NumButtons * config.NumFloors) * config.NumElevs
	bestElevator := onlineIDs[0]
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
			cost += 4
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

func sumOrders(incomming [config.NumElevs][config.NumFloors][config.NumButtons]bool) int {
	sum := 0
	for elev := 0; elev < config.NumElevs; elev++ {
		for floor := 0; floor < config.NumFloors; floor++ {
			for btn := 0; btn < config.NumButtons-1; btn++ {
				if incomming[elev][floor][btn] {
					sum++
					fmt.Println("Teller i etasje", floor)
				}
			}
		}
	}
	return sum
}
