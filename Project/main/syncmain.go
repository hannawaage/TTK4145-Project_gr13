package main

import (
	. "../config"
	. "../driver-go/elevio"
	"../sync/SyncFunc"
	 "fmt"
)

func main() {
	const NumElevs = NumElevs
	const NumFloors = NumFloors
	const NumButtons = NumButtons

	var first int = 0
	var second int = 1
	var third int = 2
	//var fourth int = 3


	var elev1 = Elevator{
		Id: 0,
		Floor: second,
		Dir: MD_Up,
		State: Moving,
		//Orders: [NumFloors][NumButtons]bool,
		//Lights: [NumElevs][NumFloors][NumButtons]bool,
	}
	elev1.Orders[first][BT_HallUp] = true
	elev1.Orders[first][BT_Cab] = true

	var elev2 = Elevator{
		Id: 1,
		Floor: third,
		Dir: MD_Down,
		State: Moving,
		//Orders: [NumFloors][NumButtons]bool,
		//Lights: [NumElevs][NumFloors][NumButtons]bool,
	}
	elev2.Orders[second][BT_HallUp] = true
	elev2.Orders[first][BT_HallUp] = true

	var elev3 = Elevator{
		Id: 2,
		Floor: first,
		Dir: MD_Up,
		State: Idle,
		//Orders: [NumFloors][NumButtons]bool,
		//Lights: [NumElevs][NumFloors][NumButtons]bool,
	}

	//elev3.Orders[third][BT_HallUp] = true
	elev3.Orders[first][BT_HallUp] = true
	elev3.Orders[second][BT_HallDown] = true

	var allOrders = [NumElevs]Elevator{elev1,elev2,elev3}
	fmt.Println("Elev1",allOrders[first].Orders)
	fmt.Println("Elev2",allOrders[second].Orders)
	fmt.Println("Elev3",allOrders[third].Orders)
	fmt.Println("AFTER FUNCTION")

	allOrders = SyncFunc.CostFunction(allOrders)
	fmt.Println("Elev1",allOrders[first].Orders)
	fmt.Println("Elev2",allOrders[second].Orders)
	fmt.Println("Elev3",allOrders[third].Orders)
}
