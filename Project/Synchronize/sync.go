package Synchronize

import (
	"../network/bcast"
	"../network/localip"
	"../network/peers"
	"../config"
)

/*
Send og motta meldinger, anta offline hvis ikke svar


*/


type SyncChns struct {
	Status chan Message
	//to be continued...
}
allOrders := make(chan)
oppdatertStates := make(chan)
oppdatertOrdre [][]bool := make(chan)
var online := make(chan bool)
kvitteringer [kvittering][id] := make(chan)



/*
Funksjon/goroutine som henter meldinger, sjekker kvitteringer,
og legger ut oppdateringer på channels som Syncro bruker til
kost funksjon/merging av ordre

Skal deles meldinger med ordre hvert sekund, svarer med kvittering,
anta offline hvis ikke to kvitteringer per melding
*/

func meldinger (
	timeout := make(chan bool)
	go func() { time.Sleep(1 * time.Second); timeout <- true }()


	select {
		case melding := <-ch.InnMelding:
			oppdatertStates = melding.Elevator
			oppdatertOrdre = melding.Ordre
			kvitteringPåMelding.sendMelding
		case kvittering := <-ch.InnKvittering:
			kvitteringer <- kvittering
		case <-timeout:
			for kvitteringer:
				to unike id per kvittering?
				if not
					online = false
		}
)



/*
online:
	master: ta inn egne og andres states og ordre fra channels, kostfunksjon i 3D,
			distribuere allOrders, legg til timespamps, sett på lys etter mottatte kvittering
	backup: send states og ordre, ta imot allOrders
!online:
	merge allOrders og nye ordre
*/
func Syncro (

	if (online) {
		if(iAmMaster) {
			heis2 := <- Elevator2chn
			heis3 := <- Elevator3chn

	liste med heiser
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
)
