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
    Elevator chan Elevator
    OrderAbove chan bool
    OrderBelow chan bool
    ShouldStop chan bool
    //to be continued...
}

timerChns := TimerChns {
    startTimer: make(chan int)
    stopTimer: make(chan bool)
    timerTimeout: make (chan bool)
}

func RunElevator(esmChns EsmChns) {
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

    */

    for {
        select {
        case buttonEvent := <- NewOrder:
          fsm_buttonRequest(buttonEvent, esmChns.Elevator.State)
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
            }*/
        }
    }
}

func fsm_buttonRequest(buttonEvent elevio.ButtonEvent, Elevator EsmChns.Elevator) {
  switch(Elevator.state) {
  case Idle:
      if(Elevator.floor == buttonEvent.floor){
        elevio.SetDoorOpenLamp(true)
        timerChn.startTimer <- config._openDoorTime
        Elevator.state = OpenDoor
      } else {

      }

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

  switch(elevator.behaviour){

   case EB_DoorOpen:
       if(elevator.floor == btn_floor){
           timer_start(elevator.config.doorOpenDuration_s);
       } else {
           elevator.requests[btn_floor][btn_type] = 1;
       }
       break;

   case EB_Moving:
       elevator.requests[btn_floor][btn_type] = 1;
       break;

   case EB_Idle:
       if(elevator.floor == btn_floor){
           outputDevice.doorLight(1);
           timer_start(elevator.config.doorOpenDuration_s);
           elevator.behaviour = EB_DoorOpen;
       } else {
           elevator.requests[btn_floor][btn_type] = 1;
           elevator.dirn = requests_chooseDirection(elevator);
           outputDevice.motorDirection(elevator.dirn);
           elevator.behaviour = EB_Moving;
       }
       break;
}
