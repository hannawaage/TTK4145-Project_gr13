package elevsm

import (
  . "../driver-go/elevio"
 //"fmt"
  . "../config"
//  . "../timer"
  . "./esmFunctions"
  "math/rand"
    "time"
)
func SyncTest(CurrentAllOrders chan<- [NumElevs][NumFloors][NumButtons]bool){
  rand.Seed(time.Now().UTC().UnixNano())
  allOrders:= [NumElevs][NumFloors][NumButtons]bool{}
    f := 0
    b := 0
    m := 0
    n := 0
    for {
      f = rand.Intn(4)
      b = rand.Intn(3)
      m = rand.Intn(4)
      n = rand.Intn(3)
      if (f == 0 && b == 1) || (f == 3 && b == 0) {
        b = 2
      }
      if (m == 0 && n == 1) || (m == 3 && n == 0) {
        n = 2
      }
      allOrders[1][f][b] = true
      //allOrders[1][m][n] = true // for å legge til 2 ordre om gangen
    	CurrentAllOrders <- allOrders
      allOrders[1][f][b] = false
      //allOrders[1][m][n] = false
      time.Sleep(3*time.Second)
    }
}


func RunElevator(esmChns EsmChns) {

    elevator := Elevator{
      State:  Idle,
      Orders: [NumFloors][NumButtons]bool{},
    }
    doorTimedOut := time.NewTimer(3 * time.Second)
    doorTimedOut.Stop()
    elevator.Floor = InitElev(elevator,esmChns)
    go ShareElev(elevator,esmChns)

    for {
      select {

      case newButtonOrder := <-esmChns.Buttons:
        if elevator.Orders[newButtonOrder.Floor][newButtonOrder.Button] == false{ //Hvis ikke allerede en ordre
          elevator.Orders[newButtonOrder.Floor][newButtonOrder.Button] = true
          go ShareElev(elevator,esmChns)
          elevator.Orders[newButtonOrder.Floor][newButtonOrder.Button] = false //Så ordren ikke påvirker esm før kostfunksjonen har evaluert den
        }

      case currentAllOrders:= <-esmChns.CurrentAllOrders:
        elevator.Orders = SetOrders(elevator, currentAllOrders)
        switch elevator.State {
        case Undefined:
        case Idle:
          elevator.Dir = SetDirection(elevator)
          SetMotorDirection(elevator.Dir)
          if elevator.Dir == MD_Stop {//if already at the correct floor
            elevator.State = DoorOpen
            SetDoorOpenLamp(true)
            doorTimedOut.Reset(3*time.Second)
            elevator.Orders = ClearOrders(elevator)
          } else {
            elevator.State = Moving
          }
        case Moving:
        case DoorOpen:
          elevator.Orders = ClearOrders(elevator)
        default:
        }
        go ShareElev(elevator,esmChns)

      case newFloor:= <-esmChns.Floors:
          elevator.Floor = newFloor
          SetFloorIndicator(newFloor)
          if ShouldStop(elevator) || (!ShouldStop(elevator) && elevator.Orders == [NumFloors][NumButtons]bool{})  {
            SetDoorOpenLamp(true)
            elevator.State = DoorOpen
            SetMotorDirection(MD_Stop)
            doorTimedOut.Reset(3 * time.Second)
            elevator.Orders = ClearOrders(elevator)
          }
          go ShareElev(elevator,esmChns)

      case <-doorTimedOut.C:
          SetDoorOpenLamp(false)
          elevator.Dir = SetDirection(elevator)
          if elevator.Dir == MD_Stop {
            elevator.State = Idle
          } else {
            elevator.State = Moving
            SetMotorDirection(elevator.Dir)
          }
          go ShareElev(elevator,esmChns)
      }
    }
}
