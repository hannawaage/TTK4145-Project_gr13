package orders

import (
	"fmt"

	. "../config"
	. "../driver-go/elevio"
	//"../StateMachine/esm"
)
// Funksjoner som begynner med stor forbokstav kan kun brukes utenfor modulen, eks UpdateOrders, motsatt for funksjoner med liten forbokstav
// . "../directory/example" gjør at man slipper example.Function, kan bare bruke Function

func UpdateState(elevator Elevator, esmChns EsmChns){
	for{
		select {
    case a := <-esmChns.Floors:
    	elevator.Floor = a
			fmt.Printf("%+v\n", elevator.Floor)
			fmt.Printf("%+v\n", elevator.Dir)
			ShouldStop(elevator)
		case b := <- esmChns.Buttons:
			SetButtonLamp(b.Button, b.Floor, true)
			elevator.Orders[b.Floor][b.Button] = true
			fmt.Printf("%+v\n", elevator.Orders)
			fmt.Printf("%+v\n", elevator.Dir)
			SetDirection(elevator)
    default:
    }
	}
}

func AddOrder(button ButtonEvent, receiver chan<- ButtonEvent) {
	receiver <- button
}

func SetDirection(elevator Elevator) {
	switch elevator.Dir {
	case MD_Up:
		if ordersAbove(elevator) {
			elevator.Dir = MD_Up
			SetMotorDirection(elevator.Dir)

		} else if ordersBelow(elevator) {
			elevator.Dir = MD_Down
			SetMotorDirection(elevator.Dir)

		}
	case MD_Down:
		if ordersBelow(elevator) {
			elevator.Dir = MD_Down
			SetMotorDirection(elevator.Dir)

		} else if ordersAbove(elevator) {
			elevator.Dir = MD_Up
			SetMotorDirection(elevator.Dir)

		}

	case MD_Stop:
		if ordersBelow(elevator) {
			elevator.Dir = MD_Down
			SetMotorDirection(elevator.Dir)

		} else if ordersAbove(elevator) {
			elevator.Dir = MD_Up
			SetMotorDirection(elevator.Dir)

		}
	}

	//elevator.Dir = MD_Stop
	SetMotorDirection(elevator.Dir)

}

func ShouldStop(elevator Elevator) {
	switch elevator.Dir {
	case MD_Up:
		if elevator.Orders[elevator.Floor][BT_HallUp] || elevator.Orders[elevator.Floor][BT_Cab] || !ordersAbove(elevator) {
			elevator.Dir = MD_Stop
			SetMotorDirection(elevator.Dir)

		}

	case MD_Down:
		if elevator.Orders[elevator.Floor][BT_HallDown] || elevator.Orders[elevator.Floor][BT_Cab] || !ordersBelow(elevator) {
			elevator.Dir = MD_Stop
			SetMotorDirection(elevator.Dir)

		}
	case MD_Stop:
	default:
		fmt.Printf("%+v\n", elevator.Orders)

	}
	//elevator.Dir = MD_Stop
	SetMotorDirection(elevator.Dir)

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


/*
func GetLocalOrders(elevator Elevator, esmChns EsmChns,receiver chan ButtonEvent) {
	for {
		buttonEvent := <-esmChns.Buttons
		SetButtonLamp(buttonEvent.Button, buttonEvent.Floor, true)
		fmt.Printf("%+v\n", elevator.Floor)
		if elevator.Floor != buttonEvent.Floor {
			elevator.Orders[buttonEvent.Floor][buttonEvent.Button] = true
		}
		fmt.Printf("%+v\n", elevator.Orders)
	}
}
*/
