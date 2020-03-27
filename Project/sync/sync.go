package sync

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"../config"
	"../network/localip"
)

func Sync(id string, syncCh config.SyncChns, esmChns config.EsmChns) {
	const numPeers = config.NumElevs - 1
	idDig, _ := strconv.Atoi(id)
	idDig--
	masterID := idDig
	var (
		elev               config.Elevator
		onlineIPs          []string
		receivedReceipt    []string
		currentMsgID       int
		numTimeouts        int
		updatedLocalOrders [config.NumElevs][config.NumFloors][config.NumButtons]bool
		currentAllOrders   [config.NumElevs][config.NumFloors][config.NumButtons]bool
		online             bool
	)
	go func() {
		for {
			select {
			case b := <-syncCh.Online:
				if b {
					online = true
					fmt.Println("Yaho, we are online!")
				} else {
					online = false
					fmt.Println("Boo, we are offline.")

				}
			case elev := <-esmChns.Elev:
				if updatedLocalOrders[idDig] != elev.Orders {
					updatedLocalOrders[idDig] = elev.Orders
				}
			}
		}
	}()

	go func() {
		for {
			if currentAllOrders != updatedLocalOrders {
				if !online {
					updatedLocalOrders = mergeAllOrders(idDig, updatedLocalOrders)
					esmChns.CurrentAllOrders <- updatedLocalOrders
				} else {
					esmChns.CurrentAllOrders <- currentAllOrders
				}
				currentAllOrders = updatedLocalOrders
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}

	msgTimer := time.NewTimer(5 * time.Second)
	msgTimer.Stop()

	go func() {
		for {
			currentMsgID = rand.Intn(256)
			msg := config.Message{elev, updatedLocalOrders, currentMsgID, false, localIP, id}
			syncCh.SendChn <- msg
			msgTimer.Reset(800 * time.Millisecond)
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		select {
		case incomming := <-syncCh.RecChn:
			recID := incomming.LocalID
			recIDDig, _ := strconv.Atoi(recID)
			recIDDig--
			if id != recID { //Hvis det ikke er fra oss selv, BYTTES TIL IP VED KJØRING PÅ FORSKJELLIGE MASKINER
				if !contains(onlineIPs, recID) {
					// Dersom heisen enda ikke er registrert, sjekker vi om vi nå er online og sjekker om vi er master
					onlineIPs = append(onlineIPs, recID)
					if len(onlineIPs) == numPeers {
						syncCh.Online <- true
						for i := 0; i < numPeers; i++ {
							theID, _ := strconv.Atoi(onlineIPs[i])
							theID--
							if theID < masterID {
								masterID = theID
							}
						}
						/*
							Dette er ved diff på IP:
							localDig, _ := strconv.Atoi(localIP[len(localIP)-3:])
							for i := 0; i <= numPeers; i++ {
								theIP := onlineIPs[i]
								lastDig, _ := strconv.Atoi(theIP[len(theIP)-3:])
								if localDig < lastDig {
									iAmMaster = false
									break
								}
							}
						*/
					}
				}

				if !incomming.Receipt {
					if online {
						if currentAllOrders[recIDDig] != incomming.AllOrders[recIDDig] {
							// Hvis vi mottar noe nytt
							if masterID == idDig {
								// Hvis jeg er master: oppdater ordrelisten vi skal sende ut med kostfunksjon
								updatedLocalOrders = CostFunction()
								//costfcn(idDig, currentAllOrders, incomming.AllOrders[recIDDig])
								esmChns.CurrentAllOrders <- currentAllOrders
								currentAllOrders = updatedLocalOrders
							} else if masterID == recIDDig {
								// Hvis meldingen er fra Master: oppdater med en gang (masters word is law)
								updatedLocalOrders = incomming.AllOrders
								esmChns.CurrentAllOrders <- currentAllOrders
								currentAllOrders = updatedLocalOrders
							}
						}
					}
					// Hvis det ikke er en kvittering, skal vi svare med kvittering
					msg := config.Message{elev, updatedLocalOrders, incomming.MsgId, true, localIP, id}
					//sender ut fem kvitteringer på femti millisekunder
					for i := 0; i < 5; i++ {
						syncCh.SendChn <- msg
						time.Sleep(10 * time.Millisecond)
					}
				} else { // Hvis det er en kvittering
					if incomming.MsgId == currentMsgID {
						if !contains(receivedReceipt, recID) {
							receivedReceipt = append(receivedReceipt, recID)
							if len(receivedReceipt) == numPeers {
								//Hvis vi har fått bekreftelse fra alle andre peers på meldingen
								numTimeouts = 0
								msgTimer.Stop()
								receivedReceipt = receivedReceipt[:0]
								// Har fått bekreftet fra resten at de har fått med seg mine nye bestillinger
							}
						}
					}
				}
			}
		case <-msgTimer.C:
			numTimeouts++
			if numTimeouts > 2 {
				syncCh.Online <- false
				fmt.Println("Three timeouts in a row")
				numTimeouts = 0
				onlineIPs = onlineIPs[:0]
				masterID = idDig
			}
		}
	}
}

func contains(elevs []string, str string) bool {
	for _, a := range elevs {
		if a == str {
			return true
		}
	}
	return false
}

func costfcn(id int, current [config.NumElevs][config.NumFloors][config.NumButtons]bool, new [config.NumFloors][config.NumButtons]bool) [config.NumElevs][config.NumFloors][config.NumButtons]bool {
	current[id] = new
	allOrderMat := mergeAllOrders(0, current)
	return allOrderMat
}

func mergeAllOrders(id int, all [config.NumElevs][config.NumFloors][config.NumButtons]bool) [config.NumElevs][config.NumFloors][config.NumButtons]bool {
	var merged [config.NumElevs][config.NumFloors][config.NumButtons]bool
	merged[id] = all[id]
	for elev := 0; elev < config.NumElevs; elev++ {
		if elev == id {
			continue
		}
		for floor := 0; floor < config.NumFloors; floor++ {
			for btn := 0; btn < config.NumButtons-1; btn++ {
				if all[elev][floor][btn] {
					merged[id][floor][btn] = true
					all[elev][floor][btn] = false
				}
			}
		}
	}
	return merged
}
