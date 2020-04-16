package main

import (
	//"flag"
	//"strconv"

	"flag"
	"strconv"

	"../config"
	. "../driver-go/elevio"
	. "../esm"
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

	var bcport string
	var id string
	flag.StringVar(&bcport, "bcport", "", "bcport of this peer")
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	idDig, _ := strconv.Atoi(id)
	idDig--

	Init(bcport, NumFloors)

	syncChns := config.SyncChns{
		SendChn:      make(chan config.Message),
		RecChn:       make(chan config.Message),
		OrderTimeout: make(chan bool),
	}

	bcastport := 16576

	go bcast.Transmitter(bcastport, syncChns.SendChn)
	go bcast.Receiver(bcastport, syncChns.RecChn)
	go sync.Sync(idDig, syncChns, esmChns)

	go PollButtons(esmChns.Buttons)
	go PollFloorSensor(esmChns.Floors)
	go RunElevator(esmChns, idDig)

	select{}
}
