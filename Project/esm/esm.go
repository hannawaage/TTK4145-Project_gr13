package esm

import (
	"time"

	. "../config"
	. "../driver-go/elevio"
)

func RunElevator(esmChns EsmChns, id int) {

	elevator := Elevator{
		Id:     id,
		State:  Idle,
		Orders: [NumFloors][NumButtons]bool{},
		Lights: [NumElevs][NumFloors][NumButtons]bool{},
	}
	doorTimedOut := time.NewTimer(DoorOpenTime)
	doorTimedOut.Stop()
	elevator.Floor = InitElev(elevator, esmChns)
	go ShareElev(elevator, esmChns)

	for {
		select {

		case newButtonOrder := <-esmChns.Buttons:
			if elevator.Orders[newButtonOrder.Floor][newButtonOrder.Button] == false { //Hvis ikke allerede en ordre
				elevator.Orders[newButtonOrder.Floor][newButtonOrder.Button] = true
				go ShareElev(elevator, esmChns)
				elevator.Orders[newButtonOrder.Floor][newButtonOrder.Button] = false //Så ordren ikke påvirker esm før kostfunksjonen har evaluert den
			}

		case currentAllOrders := <-esmChns.CurrentAllOrders:
			elevator.Orders, elevator.Lights = SetCurrentOrders(id, elevator, currentAllOrders)
			switch elevator.State {
			case Undefined:
			case Idle:
				elevator.Dir = SetDirection(elevator)
				SetMotorDirection(elevator.Dir)
				if elevator.Dir == MD_Stop {
					if OrdersInFloor(elevator) {
						elevator.State = DoorOpen
						SetDoorOpenLamp(true)
						doorTimedOut.Reset(3 * time.Second)
						elevator.Orders, elevator.Lights = ClearOrders(id, elevator)
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
			SetFloorIndicator(newFloor)
			if ShouldStop(elevator) || (!ShouldStop(elevator) && elevator.Orders == [NumFloors][NumButtons]bool{}) {
				SetDoorOpenLamp(true)
				elevator.State = DoorOpen
				SetMotorDirection(MD_Stop)
				doorTimedOut.Reset(DoorOpenTime)
				elevator.Orders, elevator.Lights = ClearOrders(id, elevator)
			}
			go ShareElev(elevator, esmChns)

		case <-doorTimedOut.C:
			SetDoorOpenLamp(false)
			elevator.Dir = SetDirection(elevator)
			if elevator.Dir == MD_Stop {
				elevator.State = Idle
			} else {
				elevator.State = Moving
				SetMotorDirection(elevator.Dir)
			}
			go ShareElev(elevator, esmChns)
		}
	}
}
