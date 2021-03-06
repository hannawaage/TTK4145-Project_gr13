package esm

import (
	"time"

	"../config"
	"../driver-go/elevio"
)

const (
	NumElevs     = config.NumElevs
	NumFloors    = config.NumFloors
	NumButtons   = config.NumButtons
	DoorOpenTime = config.DoorOpenTime
	Undefined    = config.Undefined
	Idle         = config.Idle
	Moving       = config.Moving
	DoorOpen     = config.DoorOpen
)

func RunElevator(esmChns config.EsmChns, id int) {

	elevator := config.Elevator{
		Id:     id,
		State:  config.Idle,
		Orders: [NumFloors][NumButtons]int{},
		Lights: [NumFloors][NumButtons]bool{},
	}
	doorTimedOut := time.NewTimer(DoorOpenTime)
	doorTimedOut.Stop()
	elevator.Floor = InitElev(elevator, esmChns)
	go ShareElev(elevator, esmChns)

	for {
		select {
		case newButtonOrder := <-esmChns.Buttons:
			if elevator.Orders[newButtonOrder.Floor][newButtonOrder.Button] == 0 {
				elevator.Orders[newButtonOrder.Floor][newButtonOrder.Button] = 1
				go ShareElev(elevator, esmChns)
			}

		case currentAllOrders := <-esmChns.CurrentAllOrders:
			elevator.Orders, elevator.Lights = SetCurrentOrders(id, elevator, currentAllOrders)
			switch elevator.State {
			case Undefined:
			case Idle:
				elevator.Dir = SetDirection(elevator)
				elevio.SetMotorDirection(elevator.Dir)
				if elevator.Dir == elevio.MD_Stop {
					if OrdersInFloor(elevator) {
						elevator.Orders, elevator.Lights = ClearOrders(id, elevator)
						elevator.State = DoorOpen
						elevio.SetDoorOpenLamp(true)
						doorTimedOut.Reset(3 * time.Second)
					}
				} else {
					elevator.State = Moving
				}
			case Moving:
			case DoorOpen:
				elevator.Orders, elevator.Lights = ClearOrders(id, elevator)
			default:
			}
			go ShareElev(elevator, esmChns)

		case newFloor := <-esmChns.Floors:
			elevator.Floor = newFloor
			elevio.SetFloorIndicator(newFloor)
			if ShouldStop(elevator) || (!ShouldStop(elevator) && elevator.Orders == [NumFloors][NumButtons]int{}) {
				elevio.SetDoorOpenLamp(true)
				elevator.State = DoorOpen
				elevio.SetMotorDirection(elevio.MD_Stop)
				doorTimedOut.Reset(DoorOpenTime)
				elevator.Orders, elevator.Lights = ClearOrders(id, elevator)
			}
			go ShareElev(elevator, esmChns)

		case <-doorTimedOut.C:
			elevio.SetDoorOpenLamp(false)
			elevator.Dir = SetDirection(elevator)
			if elevator.Dir == elevio.MD_Stop {
				elevator.State = Idle
			} else {
				elevator.State = Moving
				elevio.SetMotorDirection(elevator.Dir)
			}
			go ShareElev(elevator, esmChns)
		}
	}
}
