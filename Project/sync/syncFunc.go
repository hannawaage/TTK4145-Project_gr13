package sync

import (
	"../config"
	"../driver-go/elevio"
)

// CostFunction tar inn en allOrders, id, lage ny
func CostFunction(allOrders [config.NumElevs]config.Elevator, onlineIPs []int) [config.NumElevs][config.NumFloors][config.NumButtons]bool {
	var allOrdersMat [config.NumElevs][config.NumFloors][config.NumButtons]bool
	bestElevator := allOrders[0].Id
	for elevator := 0; elevator < config.NumElevs; elevator++ {
		for floor := 0; floor < config.NumFloors; floor++ {
			for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++ {
				if allOrders[elevator].Orders[floor][button] == true {
					order := elevio.ButtonEvent{
						Floor:  floor,
						Button: button,
					}
					bestElevator = costCalculator(order, &allOrders, onlineIPs)
					allOrders[elevator].Orders[order.Floor][order.Button] = false
					allOrders[bestElevator].Orders[order.Floor][order.Button] = true
				}
			}
		}
	}
	for elevator := 0; elevator < config.NumElevs; elevator++ {
		allOrdersMat[elevator] = allOrders[elevator].Orders
	}
	return allOrdersMat
}

func costCalculator(order elevio.ButtonEvent, allOrders *[config.NumElevs]config.Elevator, onlineIPs []int) int {
	minCost := (config.NumButtons * config.NumFloors) * config.NumElevs
	bestElevator := onlineIPs[0]
	for elevator := 0; elevator < config.NumElevs; elevator++ {
		if !contains(onlineIPs, allOrders[elevator].Id) {
			continue
		}
		cost := order.Floor - allOrders[elevator].Floor
		if (cost == 0) && (allOrders[elevator].State != config.Moving) {
			bestElevator = elevator
			return bestElevator
		}

		if cost < 0 {
			cost = -cost
			if allOrders[elevator].Dir == elevio.MD_Up {
				cost += 3
			}
		} else if cost > 0 {
			if allOrders[elevator].Dir == elevio.MD_Down {
				cost += 3
			}
		}
		if cost == 0 && allOrders[elevator].State == config.Moving {
			cost += 4
		}

		if allOrders[elevator].State == config.DoorOpen {
			cost++
		}

		if cost < minCost {
			minCost = cost
			bestElevator = elevator
		}
	}
	return bestElevator
}
