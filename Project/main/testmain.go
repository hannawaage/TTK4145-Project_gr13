package main

import (
	. "../Orders"
	. "../config"
	. "../driver-go/elevio"
//	"time"
//	"fmt"

)

func main() {

	numFloors := 4

	elevator := Elevator{
		Id:     1,
		Orders: [NumFloors][NumButtons]bool{},
	}

	esmChns := EsmChns{
		NewOrder: make(chan ButtonEvent),
		Buttons: make(chan ButtonEvent),
		Floors:  make(chan int),
	}
	Init("localhost:15657", numFloors)

	go PollButtons(esmChns.Buttons)
	go PollFloorSensor(esmChns.Floors)
	go UpdateState(elevator, esmChns)
	//go GetLocalOrders(elevator, esmChns, esmChns.NewOrder)
//go SetDirection(elevator)
	for {

	//	fmt.Printf("%+v\n", elevator.Orders)

		//time.Sleep(4000 * time.Millisecond)

	}
}
