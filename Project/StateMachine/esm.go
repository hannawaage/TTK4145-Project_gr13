package elevsm

import (
	"fmt"

	. "../config"
	. "../driver-go/elevio"

	//  . "../timer"
	"math/rand"
	"time"

	. "./esmFunctions"
)

func SyncTest(CurrentAllOrders chan<- [NumElevs][NumFloors][NumButtons]bool, elev chan Elevator) {
	newOrderTimer := time.NewTimer(2500 * time.Millisecond)
	newOrderTimer.Stop()
	rand.Seed(time.Now().UTC().UnixNano())
	f := 0
	b := 0
	m := 0
	n := 0
	o := 0
	p := 0
	allOrders := [NumElevs][NumFloors][NumButtons]bool{}
	for {
		select {
		case elevator := <-elev:
			for floor := 0; floor < NumFloors; floor++ {
				for btn := 0; btn < NumButtons; btn++ {
					if !elevator.Orders[floor][btn] && allOrders[1][floor][btn] { //id
						allOrders[1][floor][btn] = false // id
					}
				}
			}
			newOrderTimer.Reset(2500 * time.Millisecond)
		case <-newOrderTimer.C:
			f = rand.Intn(4)
			b = rand.Intn(3)
			m = rand.Intn(4)
			n = rand.Intn(3)
			o = rand.Intn(4)
			p = rand.Intn(3)
			if (f == 0 && b == 1) || (f == 3 && b == 0) {
				b = 2
			}
			if (m == 0 && n == 1) || (m == 3 && n == 0) {
				n = 2
			}
			if (o == 0 && p == 1) || (o == 3 && p == 0) {
				p = 2
			}
			allOrders[1][f][b] = true //random ordre til heisen
			fmt.Println("Min ordre: ", f, ",", b)
			allOrders[2][m][n] = true // random ordre fra andre heiser
			allOrders[2][o][p] = true // random ordre fra andre heiser
			CurrentAllOrders <- allOrders
			allOrders[2][m][n] = false //den andre heisen "utfører" ordren med en gang
			allOrders[2][o][p] = false //den andre heisen "utfører" ordren med en gang
		}
	}
}

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
				at := OrdersInFloor(elevator)
				fmt.Println(at)
				if elevator.Dir == MD_Stop { //&& OrdersInFloor(elevator) { //if already at the correct floor
					if at {
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
