package main

import (
	. "../StateMachine"
	. "../config"
	. "../driver-go/elevio"
//	"time"
	//"fmt"

)

func main() {

	esmChns := EsmChns{
		Elev: make(chan Elevator),
		CurrentAllOrders: make(chan [NumElevs][NumFloors][NumButtons]bool),
		Buttons: make(chan ButtonEvent),
		Floors:  make(chan int),

	}
	Init("localhost:15657", NumFloors)

	go SyncTest(esmChns.CurrentAllOrders)

	go PollButtons(esmChns.Buttons)
	go PollFloorSensor(esmChns.Floors)
	go RunElevator(esmChns)

	for {

	}
}
