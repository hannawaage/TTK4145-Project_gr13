package orders

import (
  . "../driver-go/elevio"
  "fmt"
  "../config"
  "../StateMachine"
)

type SignalChns struct {
    drv_buttons chan elevio.ButtonEvent
    drv_floors chan int
}

func UpdateOrders(esmChns EsmChns, signalChns SignalChns) {
  for {
    select {
    case buttonEvent := <-signalChns.drv_buttons:
      esmChns.NewOrder = buttonEvent
      if (esmChns.Elevator.floor != buttonEvent.Floor) {
        esmChns.Elevator.Orders[buttonEvent.Floor][buttonEvent.ButtonType] = true
      }
    }
  case floorRegistered := <-signalChns.drv_floors:
      esmChns.Elevator.Orders[floorRegistered][:] = false
    }
}

func shouldStop(elevator esmChns.Elevator) bool {
	switch elevator.Dir {
	case DirUp:
		return elevator.Orders[elevator.Floor][BT_HallUp] ||
			elevator.Orders[elevator.Floor][BT_Cab] ||
			!ordersAbove(elevator)
	case DirDown:
		return elevator.Orders[elevator.Floor][BT_HallDown] ||
			elevator.Orders[elevator.Floor][BT_Cab] ||
			!ordersBelow(elevator)
	case DirStop:
	default:
	}
	return false
}


func ordersAbove(elevator esmChns.Elevator) bool {
	for floor := elevator.Floor + 1; floor < config._numFloors; floor++ {
		for btn := 0; btn < config._numButtons; btn++ {
			if elevator.Orders[floor][btn] {
				return true
			}
		}
	}
	return false
}

func ordersBelow(elevator esmChns.Elevator) bool {
	for floor := 0; floor < elevator.Floor; floor++ {
		for btn := 0; btn < config._numButtons; btn++ {
			if elevator.Orders[floor][btn] {
				return true
			}
		}
	}
	return false
}
