package sync

import (
	"../config"
	"../driver-go/elevio"
)

// CostFunction tar inn en allElevs, id, lage ny
func CostFunction(id int, allElevs [config.NumElevs]config.Elevator, onlineIPs []int) [config.NumElevs][config.NumFloors][config.NumButtons]bool {
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
					bestElevator = costCalculator(id, order, &allElevs, onlineIPs)
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

func costCalculator(id int, order elevio.ButtonEvent, allElevs *[config.NumElevs]config.Elevator, onlineIPs []int) int {
	minCost := (config.NumButtons * config.NumFloors) * config.NumElevs
	bestElevator := onlineIPs[0]
	for elevator := 0; elevator < config.NumElevs; elevator++ {
		if !contains(onlineIPs, allElevs[elevator].Id) && (elevator != id) {
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
