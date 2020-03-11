package esmFunctions

import (
	//"fmt"
	. "../../config"
	. "../../driver-go/elevio"
)


// Funksjoner som begynner med stor forbokstav kan kun brukes utenfor modulen, eks UpdateOrders, motsatt for funksjoner med liten forbokstav
// . "../directory/example" gjør at man slipper example.Function, kan bare bruke Function




func SetDirection(elevator Elevator) MotorDirection {
	var d MotorDirection = MD_Stop
	switch elevator.Dir {
	case MD_Up:
		if ordersAbove(elevator) {
			d = MD_Up

		} else if ordersBelow(elevator) {
			d = MD_Down

		}
	case MD_Down:
		if ordersBelow(elevator) {
			d = MD_Down

		} else if ordersAbove(elevator) {
			d = MD_Up

		}

	case MD_Stop:
		if ordersBelow(elevator) {
			d = MD_Down

		} else if ordersAbove(elevator) {
			d = MD_Up

		}
	}

	return d

}

func ShouldStop(elevator Elevator) bool{
	switch elevator.Dir {
	case MD_Up:
		if elevator.Orders[elevator.Floor][BT_HallUp] || elevator.Orders[elevator.Floor][BT_Cab] || !ordersAbove(elevator) {
			return true
		}

	case MD_Down:
		if elevator.Orders[elevator.Floor][BT_HallDown] || elevator.Orders[elevator.Floor][BT_Cab] || !ordersBelow(elevator) {
			return true

		}
	case MD_Stop:
	default:
	}
	return false
}

func ordersAbove(elevator Elevator) bool {

	for floor := elevator.Floor + 1; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elevator.Orders[floor][btn] {
				return true
			}
		}
	}
	return false
}

func ordersBelow(elevator Elevator) bool {

	for floor := 0; floor < elevator.Floor; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elevator.Orders[floor][btn] {
				return true
			}
		}
	}
	return false
}
