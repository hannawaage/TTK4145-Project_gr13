package main

import (
	"flag"
	"strconv"

	"../config"
	"../driver-go/elevio"
	"../esm"
	"../network/bcast"
	"../sync"
)

func main() {
	const NumElevs = config.NumElevs
	const NumFloors = config.NumFloors
	const NumButtons = config.NumButtons

	esmChns := config.EsmChns{
		Elev:             make(chan config.Elevator),
		CurrentAllOrders: make(chan [NumElevs][NumFloors][NumButtons]bool),
		Buttons:          make(chan elevio.ButtonEvent),
		Floors:           make(chan int),
	}

	var bcport string
	var id string
	flag.StringVar(&bcport, "bcport", "", "bcport of this peer")
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	idDig, _ := strconv.Atoi(id)
	idDig--

	elevio.Init(bcport, NumFloors)

	syncChns := config.SyncChns{
		SendChn:      make(chan config.Message),
		RecChn:       make(chan config.Message),
		OrderTimeout: make(chan bool),
	}

	bcastport := 16576

	go bcast.Transmitter(bcastport, syncChns.SendChn)
	go bcast.Receiver(bcastport, syncChns.RecChn)
	go sync.Sync(idDig, syncChns, esmChns)

	go elevio.PollButtons(esmChns.Buttons)
	go elevio.PollFloorSensor(esmChns.Floors)
	go esm.RunElevator(esmChns, idDig)

	select {}
}
