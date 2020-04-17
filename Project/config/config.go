package config

import (
	"time"

	"../driver-go/elevio"
)

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
	Dir    elevio.MotorDirection
	State  ElevState
	Orders [NumFloors][NumButtons]bool
	Lights [NumElevs][NumFloors][NumButtons]bool
}

type Message struct {
	Elev      Elevator
	AllOrders [NumElevs][NumFloors][NumButtons]bool
	MsgId     int
	IsReceipt bool
	LocalIP   string
	LocalID   int
}

type EsmChns struct {
	CurrentAllOrders chan [NumElevs][NumFloors][NumButtons]bool
	Buttons          chan elevio.ButtonEvent
	Floors           chan int
	Elev             chan Elevator
	OrderAbove       chan bool
	OrderBelow       chan bool
	ShouldStop       chan bool
}

type SyncChns struct {
	SendChn      chan Message
	RecChn       chan Message
	OrderTimeout chan bool
}
