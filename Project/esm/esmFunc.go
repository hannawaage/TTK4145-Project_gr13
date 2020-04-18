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

	for elev := 0; elev < config.NumElevs; elev++ {
		for floor := 0; floor < config.NumFloors; floor++ {
			for btn = 0; btn < config.NumButtons; btn++ {
				if elev == id {
					if elevator.Lights[elev][floor][btn] != currentAllOrders[id][floor][btn] {
						elevio.SetButtonLamp(btn, floor, currentAllOrders[id][floor][btn])
						if !currentAllOrders[id][floor][btn] {
							fmt.Println("Slår av lys pga current")
						}
					}
					elevator.Orders[floor][btn] = currentAllOrders[id][floor][btn]
					elevator.Lights[elev][floor][btn] = elevator.Orders[floor][btn]
				} else {
					if btn != config.NumButtons-1 {
						if elevator.Lights[elev][floor][btn] != currentAllOrders[elev][floor][btn] {
							elevio.SetButtonLamp(btn, floor, currentAllOrders[elev][floor][btn])
							elevator.Lights[elev][floor][btn] = currentAllOrders[elev][floor][btn]
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
