package config

import (
	"time"

	. "../driver-go/elevio"
)

// Passe på at de ulike modulene ikke importerer hverandre, designe som et hieraki.

const DoorOpenTime = 3000 * time.Millisecond
const NumElevs = 3
const NumButtons = 3
const NumFloors = 4

type ElevState int

const (
	Undefined ElevState = iota - 1
	Idle
	Moving
	DoorOpen
)

type Elevator struct {
	Id     int
	Floor  int
	Dir    MotorDirection
	State  ElevState
	Orders [NumFloors][NumButtons]int
	Lights [NumFloors][NumButtons]bool
}

type Message struct {
	Elev      Elevator
	AllOrders [NumElevs][NumFloors][NumButtons]int
	MsgId     int
	IsReceipt bool
	LocalID   int
}

type EsmChns struct {
	CurrentAllOrders chan [NumElevs][NumFloors][NumButtons]int
	Buttons          chan ButtonEvent
	Floors           chan int
	Elev             chan Elevator
	OrderAbove       chan bool
	OrderBelow       chan bool
	ShouldStop       chan bool
}

type SyncChns struct {
	SendChn chan Message
	RecChn  chan Message
}
