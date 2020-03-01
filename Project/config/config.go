package config

import (
	"time"

	"../driver-go/elevio"
)

const OpenDoorTime = 3000 * time.Millisecond
const NumElevs = 3
const NumButtons = 3
const NumFloors = 4

type Elevator struct {
	ID          int //eller noe for å vit om master eller ikke
	Floor       int
	Dir         elevio.MotorDirection
	State       esm.ElevState
	localOrders [NumFloors][NumButtons]bool
}

type BackupMessage struct {
	Elev      Elevator
  allOrders [NumElevs][NumFloors][NumButtons]bool
}

type MasterMessage struct {
	Elev      Elevator
  allOrders [NumElevs][NumFloors][NumButtons]bool
}

