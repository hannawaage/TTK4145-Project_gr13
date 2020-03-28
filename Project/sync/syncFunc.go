package sync

import (
	"fmt"

	"../config"
	"../driver-go/elevio"
)

// CostFunction tar inn en allOrders, id, lage ny
func CostFunction(allOrders [config.NumElevs]config.Elevator) [config.NumElevs][config.NumFloors][config.NumButtons]bool {
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
					// finding best elevator for that order
					bestElevator = costCalculator(order, &allOrders)
					fmt.Println(bestElevator)
					allOrders[bestElevator].Orders[order.Floor][order.Button] = true
				}
			}
		}
	}
	for elevator := 0; elevator < config.NumElevs; elevator++ {
		allOrdersMat[elevator] = allOrders[elevator].Orders
	}
	if allOrdersMat[2][0][0] {
		fmt.Println("Etter kost har fortsatt heis 3 en ordre")
	}
	return allOrdersMat
}

func costCalculator(order elevio.ButtonEvent, allOrders *[config.NumElevs]config.Elevator) int {
	//ButtonInside??
	minCost := (config.NumButtons * config.NumFloors) * config.NumElevs
	bestElevator := allOrders[0].Id
	for elevator := 0; elevator < config.NumElevs; elevator++ {

		cost := order.Floor - allOrders[elevator].Floor
		// delete the order
		allOrders[elevator].Orders[order.Floor][order.Button] = false

		if cost == 0 && allOrders[elevator].State != config.Moving {
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
