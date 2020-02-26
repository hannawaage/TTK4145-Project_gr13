package elevsm

import (
  . "../driver-go/elevio"
  "fmt"
  "../main"
  "../config"
)


type ElevState struct {
  Idle = iota
  OpenDoor
  Moving
}

type EsmChns struct {
    NewOrder chan elevio.ButtonEvent
    Elevator chan config.Elevator
    OrderAbove chan bool
    OrderBelow chan bool
    ShouldStop chan bool
    SignalChns chan orders.SignalChns
    //to be continued...
}

timerChns := TimerChns {
    startTimer: make(chan int)
    stopTimer: make(chan bool)
    timerTimeout: make (chan bool)
}
/*
func RunElevator(esmChns EsmChns) {
    elevator := esmChns.Elevator{
      //ID:  int //eller noe for å vit om master eller ikke
      Floor: int
      Dir: elevio.MotorDirection
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
            elevio.SetButtonLamp(a.Button, a.Floor, true)
            dest = a.Floor
            if dest <  last_floor {
                d = elevio.MD_Down
            } else if dest > last_floor{
                d = elevio.MD_Up
            }
            elevio.SetMotorDirection(d)

        case a := <- signal_channel.drv_floors:
          last_floor = a
          if dest == a {
            d = elevio.MD_Stop
            elevio.SetButtonLamp(0, a, false)
            elevio.SetButtonLamp(1, a, false)
            elevio.SetButtonLamp(2, a, false)
          }
          elevio.SetMotorDirection(d)*/

/*  if dest <  a {
      d = elevio.MD_Down
  } else if dest > a{
      d = elevio.MD_Up
  } else

        case a := <- drv_obstr:
            fmt.Printf("%+v\n", a)
            if a {
                elevio.SetMotorDirection(elevio.MD_Stop)
            } else {
                elevio.SetMotorDirection(d)
            }

        case a := <- drv_stop:
            fmt.Printf("%+v\n", a)
            for f := 0; f < numFloors; f++ {
                for b := elevio.ButtonType(0); b < 3; b++ {
                    elevio.SetButtonLamp(b, f, false)
                }
            }//
        }
    }
}
*/

func fsm_buttonRequest(buttonEvent elevio.ButtonEvent, esmChns EsmChns) {
  switch(esmChns.Elevator.State) {
  case Idle:
      if esmChns.Elevator.Floor == buttonEvent.Floor{
        elevio.SetDoorOpenLamp(true)
        timerChn.startTimer <- OpenDoorTime
        esmChns.Elevator.State = OpenDoor
      } else {
        esmChns.Elevator.Orders[buttonEvent.Floor][buttonEvent.Button] = true
        esmChns.Elevator.Dir = fsm_chooseDirection(esmChns.Elevator)
        elev_io.SetMotorDirection(esmChns.Elevator.Dir)
        esmChns.Elevator.State = Moving
      }
  case OpenDoor:
    if esmChns.Elevator.Floor =  buttonEvent.floor{
      timerChn.startTimer <- OpenDoorTime
    } else {
      esmChns.Elevator.Orders[buttonEvent.Floor][buttonEvent.Button] = true
    }
  case Moving:
    esmChns.Elevator.Orders[buttonEvent.Floor][buttonEvent.Button] = true

  // sette lys??

  }



func fsm_chooseDirection(elevator EsmChns.Elevator) {
  switch elevator.Dir {
	case DirStop:
		if ordersAbove(elevator) {
			return DirUp
		} else if ordersBelow(elevator) {
			return DirDown
		} else {
			return DirStop
		}
	case DirUp:
		if ordersAbove(elevator) {
			return DirUp
		} else if ordersBelow(elevator) {
			return DirDown
		} else {
			return DirStop
		}

	case DirDown:
		if ordersBelow(elevator) {
			return DirDown
		} else if ordersAbove(elevator) {
			return DirUp
		} else {
			return DirStop
		}
	}
	return DirStop
}
