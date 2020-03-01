package Synchronize

import (
	"../network/bcast"
	"../network/localip"
	"../network/peers"
	"../config"
)

type SyncChns struct {
	Status chan Message
	//to be continued...
}


timeout := make(chan bool)



/* - les inn om alle er på nett
- kostfunksjon i 3D i så tilfelle, og distribuere (bare master pc)
- hvis ikke, merge ordre (hver pc uten nett)


//Hvordan sjekke online status????

*/

if (online) {
	if(iAmMaster) {
		heis2 := <- Elevator2chn
		heis3 := <- Elevator3chn

		kostfunksjon(esmchns.Elevator, heis2, heis3) => allOrders
		channel <- allOrders
		legg til time stamps
		sett på lys hos alle

	}
	else {
		state_channel <- esmchns.Elevator.State
		order_channel <- esmchns.Elevator.Orders
		backup := <- BackupChannel (mottar oppdatert ordreliste)
		esmchns.Elevator.Orders <- backup(ID)
	}
}
else {
	update backup with localOrders
	esmchns.Elevator.Orders <-  merge(backup)
}




