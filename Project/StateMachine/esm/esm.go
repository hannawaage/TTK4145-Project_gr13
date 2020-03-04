package elevsm

import (
  . "../../driver-go/elevio"
  _ "fmt"
  . "../../config"
  . "../../timer"
  . "../../orders"
)


/*
func RunElevator(esmChns EsmChns) {
    elevator := esmChns.Elevator{
      //ID:  int //eller noe for å vit om master eller ikke
      Floor: int
      Dir: MotorDirection
      State:  esm.ElevState
      Orders: [NumFloors][NumButtons]bool
    }
    esmChns.Elevator.State = Idle
    /*
    Alt dette i annen modul:
      if (connected)
        send iAmAlive
        if(iAmMaster)
          motta ordre
          kjør kostfunksjon
          distribuer ordre
        else
          send lokale ordre
          motta ferdig ordreliste
     else
        ordre = lokale ordre
   til slutt, uansett hvor ordre kommer fra:
   avgjør retning basert på ordre
//


    for {
        select {
        case buttonEvent := <- NewOrder:
          fsm_buttonRequest(buttonEvent, esmChns)
       /*
        case a := <- signal_channel.drv_buttons:
            fmt.Printf("%+v\n", a)
            SetButtonLamp(a.Button, a.Floor, true)
            dest = a.Floor
            if dest <  last_floor {
                d = MD_Down
            } else if dest > last_floor{
                d = MD_Up
            }
            SetMotorDirection(d)

        case a := <- signal_channel.drv_floors:
          last_floor = a
          if dest == a {
            d = MD_Stop
            SetButtonLamp(0, a, false)
            SetButtonLamp(1, a, false)
            SetButtonLamp(2, a, false)
          }
          SetMotorDirection(d)*/

/*  if dest <  a {
      d = MD_Down
  } else if dest > a{
      d = MD_Up
  } else

        case a := <- drv_obstr:
            fmt.Printf("%+v\n", a)
            if a {
                SetMotorDirection(MD_Stop)
            } else {
                SetMotorDirection(d)
            }

        case a := <- drv_stop:
            fmt.Printf("%+v\n", a)
            for f := 0; f < numFloors; f++ {
                for b := ButtonType(0); b < 3; b++ {
                    SetButtonLamp(b, f, false)
                }
            }//
        }
    }
}
*/

func esm_ButtonRequest(
  buttonEvent ButtonEvent,
  elevator Elevator,
  timerChns TimerChns) {

    switch(elevator.State) {
    case Idle:
        if elevator.Floor == buttonEvent.Floor{
          SetDoorOpenLamp(true)
          //timerChns.StartTimer <- DoorOpenTime
          elevator.State = DoorOpen
        } else {
          elevator.Orders[buttonEvent.Floor][buttonEvent.Button] = true
          elevator.Dir = ChooseDirection(elevator)
          SetMotorDirection(elevator.Dir)
          elevator.State = Moving
        }
    case DoorOpen:
      if elevator.Floor == buttonEvent.Floor{
        //timerChns.StartTimer <- DoorOpenTime
      } else {
        elevator.Orders[buttonEvent.Floor][buttonEvent.Button] = true
      }
    case Moving:
      elevator.Orders[buttonEvent.Floor][buttonEvent.Button] = true

    // sette lys??

  }
}
