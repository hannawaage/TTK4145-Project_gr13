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

func SetCurrentOrders(id int, elevator config.Elevator, currentAllOrders [config.NumElevs][config.NumFloors][config.NumButtons]bool) [config.NumFloors][config.NumButtons]bool {
	var btn elevio.ButtonType
	fmt.Println("Lokal liste er: ")
	fmt.Println(currentAllOrders[id])
	for elev := 0; elev < config.NumElevs; elev++ {
		for floor := 0; floor < config.NumFloors; floor++ {
			for btn = 0; btn < config.NumButtons; btn++ {
				if elev == id {
					if currentAllOrders[elev][floor][btn] {
						elevio.SetButtonLamp(btn, floor, true)
						elevator.Orders[floor][btn] = true
					} else {
						elevio.SetButtonLamp(btn, floor, false)
						elevator.Orders[floor][btn] = false
					}
				} else {
					if currentAllOrders[elev][floor][btn] && (btn != elevio.BT_Cab) {
						elevio.SetButtonLamp(btn, floor, true)
					} else {
						elevio.SetButtonLamp(btn, floor, false)
					}
				}

			}
		}
	}
	return elevator.Orders
}

func ClearOrders(id int, elevator config.Elevator) [config.NumFloors][config.NumButtons]bool {
	var b elevio.ButtonType
	for b = 0; b < config.NumButtons; b++ {
		elevio.SetButtonLamp(b, elevator.Floor, false)
		elevator.Orders[elevator.Floor][b] = false
	}
	return elevator.Orders
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
