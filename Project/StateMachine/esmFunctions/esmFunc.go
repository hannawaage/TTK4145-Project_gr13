package esmFunctions

import (
	"fmt"

	. "../../config"
	. "../../driver-go/elevio"
)

// Funksjoner som begynner med stor forbokstav kan kun brukes utenfor modulen, eks UpdateOrders, motsatt for funksjoner med liten forbokstav
// . "../directory/example" gjør at man slipper example.Function, kan bare bruke Function

func InitElev(elevator Elevator, esmChns EsmChns) int {
	SetMotorDirection(MD_Down)
	a := <-esmChns.Floors
	for a == -1 {
		a = <-esmChns.Floors
	}
	SetMotorDirection(MD_Stop)
	SetFloorIndicator(a)
	fmt.Println("Heisen er klar i etasje nr ", a)
	return a
}

func ShareElev(elevator Elevator, esmChns EsmChns) {
	esmChns.Elev <- elevator
}

func SetOrders(idDig int, elevator Elevator, currentAllOrders [NumElevs][NumFloors][NumButtons]bool) [NumFloors][NumButtons]bool {
	var btn ButtonType
	for elev := 0; elev < NumElevs; elev++ {
		for floor := 0; floor < NumFloors; floor++ {
			for btn = 0; btn < NumButtons; btn++ {
				if !currentAllOrders[elev][floor][btn] && elevator.Lights[elev][floor][btn] {
					elevator.Lights[elev][floor][btn] = false
					if !elevator.Orders[floor][btn] {
						SetButtonLamp(btn, floor, false)
					}
				}
				if currentAllOrders[elev][floor][btn] && !(elev != 1 && btn == NumButtons-1) { //id, hvis det ikke er cab hos annen heis
					elevator.Lights[elev][floor][btn] = true
					SetButtonLamp(btn, floor, true)
					if elev == idDig { // id
						elevator.Orders[floor][btn] = true
						fmt.Println("Updated order for elevator ", idDig+1)
					}
				}
			}
		}
	}
	return elevator.Orders, elevator.Lights
}

func ClearOrders(elevator Elevator) ([NumFloors][NumButtons]bool, [NumElevs][NumFloors][NumButtons]bool) {
	var b ButtonType
	for b = 0; b < NumButtons; b++ {
		elevator.Lights[1][elevator.Floor][b] = false //id
		SetButtonLamp(b, elevator.Floor, false)
		elevator.Orders[elevator.Floor][b] = false
	}
	return elevator.Orders, elevator.Lights
}

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

func ShouldStop(elevator Elevator) bool {
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
