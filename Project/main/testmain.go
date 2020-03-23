package main

import (
	"flag"

	. "../StateMachine"
	. "../Synchronize"
	. "../config"
	. "../driver-go/elevio"

	"../network/bcast"
	//	"time"
	//"fmt"
)

func main() {

	esmChns := EsmChns{
		Elev:             make(chan Elevator),
		CurrentAllOrders: make(chan [NumElevs][NumFloors][NumButtons]bool),
		Buttons:          make(chan ButtonEvent),
		Floors:           make(chan int),
	}
	Init("localhost:15657", NumFloors)

	/////// DETTE ER FRA SYNC ////////////
	syncChns := sync.SyncChns{
		SendChn:   make(chan sync.Message),
		RecChn:    make(chan sync.Message),
		Online:    make(chan bool),
		IAmMaster: make(chan bool),
	}

	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	port := 16576

	go bcast.Transmitter(port, syncChns.SendChn)
	go bcast.Receiver(port, syncChns.RecChn)
	go sync.Sync(id, syncChns, esmChns)
	go sync.OrdersDist(syncChns)
	/////////////

	//go SyncTest(esmChns.CurrentAllOrders)

	go PollButtons(esmChns.Buttons)
	go PollFloorSensor(esmChns.Floors)
	go RunElevator(esmChns)

	for {

	}
}
