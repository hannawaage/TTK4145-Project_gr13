package main

import (
	. "../Orders"
	. "../config"
	. "../driver-go/elevio"
)

func main() {

	numFloors := 4

	elevator := Elevator{
		Id:     1,
		Floor:  0,
		Dir:    MD_Stop,
		State:  Idle,
		Orders: [NumFloors][NumButtons]bool{},
	}

	esmChns := EsmChns{
		NewOrder: make(chan ButtonEvent),
	}

	signalChns := SignalChns{
		Buttons: make(chan ButtonEvent),
		Floors:  make(chan int),
	}

	Init("localhost:15657", numFloors)

	var d MotorDirection = MD_Up
	SetMotorDirection(d)

	/*drv_buttons := make(chan ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)*/

	go PollButtons(signalChns.Buttons)
	go PollFloorSensor(signalChns.Floors)
	go GetLocalOrders(elevator, esmChns, signalChns)
	go ShouldStop(elevator)
	go SetDirection(elevator)
	//go PollObstructionSwitch(drv_obstr)
	//go PollStopButton(drv_stop)

	for {

		//fmt.Printf("%+v\n", elevator.Orders)

		//time.Sleep(4000 * time.Millisecond)

	}
}
