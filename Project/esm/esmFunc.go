package esm

import (
	"fmt"

	"../config"
	"../driver-go/elevio"
)

func InitElev(elevator config.Elevator, esmChns config.EsmChns) int {
	elevio.SetMotorDirection(elevio.MD_Down)
	floor := <-esmChns.Floors
	for floor == -1 {
		floor = <-esmChns.Floors
	}
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetFloorIndicator(floor)
	fmt.Println("Heisen er klar i etasje nr ", floor)
	return floor
}

func ShareElev(elevator config.Elevator, esmChns config.EsmChns) {
	esmChns.Elev <- elevator
}

func SetCurrentOrders(id int, elevator config.Elevator, currentAllOrders [config.NumElevs][config.NumFloors][config.NumButtons]bool) ([config.NumFloors][config.NumButtons]bool, [config.NumElevs][config.NumFloors][config.NumButtons]bool) {
	var btn elevio.ButtonType

	/*
		Hvis det ikke er ordre i currentAll og lyset er på, skal lyset slås av
	*/

	for elev := 0; elev < config.NumElevs; elev++ {
		for floor := 0; floor < config.NumFloors; floor++ {
			for btn = 0; btn < config.NumButtons; btn++ {
				if !currentAllOrders[elev][floor][btn] && elevator.Lights[elev][floor][btn] {
					elevator.Lights[elev][floor][btn] = false
					elevio.SetButtonLamp(btn, floor, false)
					/*if !elevator.Orders[floor][btn] {

					}*/
				}
				if currentAllOrders[elev][floor][btn] && !(elev != id && btn == config.NumButtons-1) {
					elevator.Lights[elev][floor][btn] = true
					elevio.SetButtonLamp(btn, floor, true)
				}
				if elev == id {
					if currentAllOrders[id][floor][btn] {
						elevator.Orders[floor][btn] = true
					} else {
						if elevator.Orders[floor][btn] {
							elevator.Orders[floor][btn] = false
						}
					}
				}
			}
		}
	}
	return elevator.Orders, elevator.Lights
}

func ClearOrders(id int, elevator config.Elevator) ([config.NumFloors][config.NumButtons]bool, [config.NumElevs][config.NumFloors][config.NumButtons]bool) {
	var btn elevio.ButtonType
	for btn = 0; btn < config.NumButtons; btn++ {
		elevator.Lights[id][elevator.Floor][btn] = false
		elevio.SetButtonLamp(btn, elevator.Floor, false)
		elevator.Orders[elevator.Floor][btn] = false
	}
	return elevator.Orders, elevator.Lights
}

func SetDirection(elevator config.Elevator) elevio.MotorDirection {
	var d elevio.MotorDirection = elevio.MD_Stop
	switch elevator.Dir {
	case elevio.MD_Up:
		if ordersAbove(elevator) {
			d = elevio.MD_Up
		} else if ordersBelow(elevator) {
			d = elevio.MD_Down
		}
	case elevio.MD_Down:
		if ordersBelow(elevator) {
			d = elevio.MD_Down
		} else if ordersAbove(elevator) {
			d = elevio.MD_Up
		}
	case elevio.MD_Stop:
		if ordersBelow(elevator) {
			d = elevio.MD_Down
		} else if ordersAbove(elevator) {
			d = elevio.MD_Up
		}
	}
	return d
}

func ShouldStop(elevator config.Elevator) bool {
	switch elevator.Dir {
	case elevio.MD_Up:
		if elevator.Orders[elevator.Floor][elevio.BT_HallUp] || elevator.Orders[elevator.Floor][elevio.BT_Cab] || !ordersAbove(elevator) {
			return true
		}
	case elevio.MD_Down:
		if elevator.Orders[elevator.Floor][elevio.BT_HallDown] || elevator.Orders[elevator.Floor][elevio.BT_Cab] || !ordersBelow(elevator) {
			return true
		}
	case elevio.MD_Stop:
	default:
	}
	return false
}

func ordersAbove(elevator config.Elevator) bool {
	for floor := elevator.Floor + 1; floor < config.NumFloors; floor++ {
		for btn := 0; btn < config.NumButtons; btn++ {
			if elevator.Orders[floor][btn] {
				return true
			}
		}
	}
	return false
}

func ordersBelow(elevator config.Elevator) bool {
	for floor := 0; floor < elevator.Floor; floor++ {
		for btn := 0; btn < config.NumButtons; btn++ {
			if elevator.Orders[floor][btn] {
				return true
			}
		}
	}
	return false
}

func OrdersInFloor(elevator config.Elevator) bool {
	for btn := 0; btn < config.NumButtons; btn++ {
		if elevator.Orders[elevator.Floor][btn] {
			return true
		}
	}
	return false
}
