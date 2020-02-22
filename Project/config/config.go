package config

import "./driver-go/elevio"

const _openDoorTime = 3000 * time.Millisecond

type Elevator struct {
  ID int //eller noe for å vit om master eller ikke
  Floor int
  Dir elevio.MotorDirection
  State esm.ElevState
  Orders [_numFloors][_numButtons]bool
}
