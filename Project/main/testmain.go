package main

import (
	//. "../StateMachine/esmFunctions"
	. "../StateMachine"
	. "../config"
	. "../driver-go/elevio"
//	"time"
	//"fmt"

)

/*func GetOrders(button chan ButtonEvent, receiver chan ButtonEvent){
	fmt.Printf("Heihei\n")

	for{
		select{
		case buttonEvent := <-button:
			SetButtonLamp(buttonEvent.Button,buttonEvent.Floor, true)
			receiver <- buttonEvent
		default:
		}
	}
}

/*func addOrder(receiver chan<- ButtonEvent, button ButtonEvent){
	receiver <- button
}*/
func main() {


	esmChns := EsmChns{
		Elev: make(chan Elevator),
		NewOrder: make(chan ButtonEvent),
		Buttons: make(chan ButtonEvent),
		Floors:  make(chan int),

	}
	Init("localhost:15657", NumFloors)

	go PollButtons(esmChns.Buttons)
	go PollFloorSensor(esmChns.Floors)
	RunElevator(esmChns)
	//go GetOrders(esmChns.Buttons, esmChns.NewOrder)
	//go addOrder(esmChnns)

	for {

	}
}
