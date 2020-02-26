package config

import "./driver-go/elevio"

const OpenDoorTime = 3000 * time.Millisecond
const NumButtons = 3
const NumFloors = 4

type Elevator struct {
  ID int //eller noe for å vit om master eller ikke
  Floor int
  Dir elevio.MotorDirection
  State esm.ElevState
  Orders [NumFloors][NumButtons]bool
}
