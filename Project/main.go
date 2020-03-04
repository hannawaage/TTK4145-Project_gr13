package main

imporg (
  "./driver-go/elevio"
  "./StateMachine"
  "fmt"
  "./orders"
)

func main(){

    //numFloors := 4

    //Finne master og backups, lagre ID
    Init, finn peers:
    - send udp "Hei jeg lever"
    - lytt etter udp melding, legg til peers som svarer
    - lagre ID = ip adresse 5 siffer fra peers
    - sammenlign siste siffer ID:
      hvis lavest -> master
      ellers -> backup 
    

    elevio.Init("localhost:15657", NumFloors)

    var d elevio.MotorDirection = elevio.MD_Stop
    elevio.SetMotorDirection(d)
    var dest int
    var last_floor int
    //lamps := make([]elevio.ButtonType, 1)

    OrderSignalChns := orders.SignalChns {
      drv_buttons:    make(chan elevio.ButtonEvent),
      drv_floors:    make(chan int)
    }


    go elevio.PollButtons(drv_buttons/*, drv_destination*/)
    go elevio.PollFloorSensor(drv_floors)
    //go elevio.PollObstructionSwitch(drv_obstr)
    //go elevio.PollStopButton(drv_stop)
    go orders.UpdateOrders(OrderSignalChns)
    go esm.RunElevator()
}
