package elevsm

import (
  . "../driver-go/elevio"
 "fmt"
  . "../config"
//  . "../timer"
  . "./esmFunctions"
)

func RunElevator(esmChns EsmChns) {
    fmt.Printf("Hello\n")

    elevator := Elevator{
      Floor: GetFloor(),
      Dir: 0,
      State:  Idle,
      Orders: [NumFloors][NumButtons]bool{},
    }
    //doorTimedOut := time.NewTimer(3 * time.Second)
    //init()
  //  esmChns.Elev <- elevator
    fmt.Printf("Heihei\n")
    fmt.Printf("%+v\n", elevator.Floor)

    for {
      select {
      case newOrder:= <-esmChns.Buttons:
        fmt.Printf("tittei\n")

          SetButtonLamp(newOrder.Button,newOrder.Floor, true)
          elevator.Orders[newOrder.Floor][newOrder.Button] = true
          fmt.Printf("%+v\n",elevator.Orders)

          switch elevator.State {
          case Undefined:
          case Idle:
            elevator.Dir = SetDirection(elevator)
            SetMotorDirection(elevator.Dir)
            if elevator.Dir == MD_Stop {
              elevator.State = DoorOpen
              //SetDoorOpenLamp()
              //SetDoorOpeTimer
              elevator.Orders[newOrder.Floor][newOrder.Button] = false
            } else {
              elevator.State = Moving
            }
          case Moving:
          case DoorOpen:
            if elevator.Floor == newOrder.Floor {
              elevator.Orders[newOrder.Floor][newOrder.Button] = false
            }
          default:
          }
          esmChns.Elev <- elevator

      case floor := <-esmChns.Floors:
        fmt.Printf("heeeello\n")
        fmt.Printf("%+v\n", elevator.Floor)


          elevator.Floor = floor
          if ShouldStop(elevator) || (!ShouldStop(elevator) &&
          elevator.Orders == [NumFloors][NumButtons]bool{})  {
            SetMotorDirection(MD_Stop)
            //SetDoorOpenLamp
            //SetDoorOpenTimer
            elevator.State = DoorOpen
          }
          esmChns.Elev <- elevator

      /*case <-doorTimeOut:
          // SetDoorOpenLamp off
          elevator.Dir = SetDirection(elevator)
          if elevator.Dir == MD_Stop {
            elevator.State = Idle
          } else {
            elevator.State = Moving
            SetMotorDirection(elevator.Dir)
          }
          esmChns.Elev <- elevator*/
      }
    }
}
