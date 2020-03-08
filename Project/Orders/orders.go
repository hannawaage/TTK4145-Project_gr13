package orders

import (
  . "../driver-go/elevio"
  //"fmt"
  . "../config"
  //"../StateMachine/esm"
)

type SignalChns struct {
    drv_buttons chan ButtonEvent
    drv_floors chan int
}

// Funksjoner som begynner med stor forbokstav kan kun brukes utenfor modulen, eks UpdateOrders, motsatt for funksjoner med liten forbokstav
// . "../directory/example" gjør at man slipper example.Function, kan bare bruke Function


/*
func UpdateOrders(elevator Elevator, esmChns chan<- EsmChns,signalChns chan<-  SignalChns) {
  for {
    select {
    case buttonEvent := <-signalChns.drv_buttons:
      esmChns.NewOrder = buttonEvent
      if (elevator.Floor != buttonEvent.Floor) {
        elevator.Orders[buttonEvent.Floor][buttonEvent.Button] = true
      }
    case floorRegistered := <-signalChns.drv_floors:
      elevator.Orders[floorRegistered][:] = false
    }
  }
}
*/



func ChooseDirection(elevator Elevator) MotorDirection {
	switch elevator.Dir {
	case MD_Up:
		if ordersAbove(elevator) {
			return MD_Up
		} else if ordersBelow(elevator) {
			return MD_Down
		}
	case MD_Down:
    if ordersBelow(elevator) {
			return MD_Down
		} else if ordersAbove(elevator) {
			return MD_Up
		}

	case MD_Stop:
		if ordersBelow(elevator) {
			return MD_Down
		} else if ordersAbove(elevator) {
			return MD_Up
		}
	}
	return MD_Stop
}


func ShouldStop(elevator Elevator) bool {
	switch elevator.Dir {
	case MD_Up:
		return elevator.Orders[elevator.Floor][BT_HallUp] || elevator.Orders[elevator.Floor][BT_Cab] || !ordersAbove(elevator)
	case MD_Down:
		return elevator.Orders[elevator.Floor][BT_HallDown] || elevator.Orders[elevator.Floor][BT_Cab] || !ordersBelow(elevator)
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
