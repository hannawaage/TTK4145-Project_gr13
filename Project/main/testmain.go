package main

import (
	//"flag"
	//"strconv"

	"flag"
	"strconv"

	. "../StateMachine"
	"../config"
	. "../driver-go/elevio"
	"../network/bcast"
	"../sync"
	//	"time"
	//"fmt"
)

func main() {
	const NumElevs = config.NumElevs
	const NumFloors = config.NumFloors
	const NumButtons = config.NumButtons
	esmChns := config.EsmChns{
		Elev:             make(chan config.Elevator),
		CurrentAllOrders: make(chan [NumElevs][NumFloors][NumButtons]bool),
		Buttons:          make(chan ButtonEvent),
		Floors:           make(chan int),
	}

	Init("localhost:12347", NumFloors)

	/////// DETTE ER FRA SYNC ////////////
	syncChns := config.SyncChns{
		SendChn:   make(chan config.Message),
		RecChn:    make(chan config.Message),
		Online:    make(chan bool),
		IAmMaster: make(chan bool),
		OfflineUpdate make(chan [NumElevs][NumFloors][NumButtons]bool),
	}

	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()
	idDig, _ := strconv.Atoi(id)
	idDig--

	bcastport := 16576

	go bcast.Transmitter(bcastport, syncChns.SendChn)
	go bcast.Receiver(bcastport, syncChns.RecChn)
	go sync.Sync(id, syncChns, esmChns)
	go sync.OrdersDistribute(idDig, syncChns, esmChns)
	/////////////

	//go SyncTest(esmChns.CurrentAllOrders, esmChns.Elev)

	go PollButtons(esmChns.Buttons)
	go PollFloorSensor(esmChns.Floors)
	go RunElevator(esmChns, idDig)

	for {

	}
}
